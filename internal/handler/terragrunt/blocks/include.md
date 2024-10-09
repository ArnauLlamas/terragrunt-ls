include
---
The `include` block is used to specify inheritance of Terragrunt configuration files. The included config (also called the `parent`) will be merged with the current configuration (also called the child) before processing. You can learn more about the inheritance properties of Terragrunt in the Filling in remote state settings with Terragrunt section of the “Keep your remote state configuration DRY” use case overview.

You can have more than one `include` block, but each one must have a unique label. It is recommended to always label your `include` blocks. Bare includes (`include` block with no label - e.g., `include` {}) are currently supported for backward compatibility, but is deprecated usage and support may be removed in the future.

`include` blocks support the following arguments:
* `name` (label): You can define multiple `include` blocks in a single terragrunt config. Each `include` block must be labeled with a unique name to differentiate it from the other includes. E.g., if you had a block `include "remote" {}`, you can reference the relevant exposed data with the expression `include.remote`.
* `path` (attribute): Specifies the path to a Terragrunt configuration file (the `parent` config) that should be merged with this configuration (the `child` config).
* `expose` (attribute, optional): Specifies whether or not the included config should be parsed and exposed as a variable. When `true`, you can reference the data of the included config under the variable `include`. Defaults to `false`. Note that the `include` variable is a map of `include` labels to the parsed configuration value.
* `merge_strategy` (attribute, optional): Specifies how the included config should be merged. Valid values are: `no_merge` (do not merge the included config), `shallow` (do a shallow merge - default), `deep` (do a deep merge of the included config).

**NOTE**: At this time, Terragrunt only supports a single level of `include` blocks. That is, Terragrunt will error out if an included config also has an `include` block defined. If you are interested in this feature, please follow https://github.com/gruntwork-io/terragrunt/issues/1566 to be notified when nested `include` blocks are supported.

Special case for shallow merge: When performing a shallow merge, all attributes and blocks are merged shallowly with replacement, except for `dependencies` blocks (NOT `dependency` block). `dependencies` blocks are deep merged: that is, all the lists of paths from included configurations are concatenated together, rather than replaced in override fashion.

Examples:

*Single include*

```hcl
# If you have the following folder structure, and the following contents for ./child/terragrunt.hcl, this will `include`
# and merge the items in the terragrunt.hcl file at the root.
#
# .
# ├── terragrunt.hcl
# └── child
#     └── terragrunt.hcl
include "root" {
  path   = find_in_parent_folders()
  expose = true
}

inputs = {
  remote_state_config = include.root.remote_state
}
```

*Multiple includes*

```hcl
# If you have the following folder structure, and the following contents for ./child/terragrunt.hcl, this will `include`
# and merge the items in the terragrunt.hcl file at the root, while only loading the data in the region.hcl
# configuration.
#
# .
# ├── terragrunt.hcl
# ├── region.hcl
# └── child
#     └── terragrunt.hcl
include "remote_state" {
  path   = find_in_parent_folders()
  expose = true
}

include "region" {
  path           = find_in_parent_folders("region.hcl")
  expose         = true
  merge_strategy = "no_merge"
}

inputs = {
  remote_state_config = include.remote_state.remote_state
  region              = include.region.region
}
```

**Limitations on accessing exposed config**

In general, you can access all attributes on `include` when they are exposed (e.g., `include.locals`, `include.inputs`, etc).

However, to support `run-all`, Terragrunt is unable to expose all attributes when the included config has a `dependency` block. To understand this, consider the following example:

*root terragrunt.hcl*

```hcl
dependency "vpc" {
  config_path = "${get_terragrunt_dir()}/../vpc"
}

inputs = {
  vpc_name = dependency.vpc.outputs.name
}
```

*child terragrunt.hcl*

```hcl
include "root" {
  path   = find_in_parent_folders()
  expose = true
}

dependency "alb" {
  config_path = (
    include.root.inputs.vpc_name == "mgmt"
    ? "../alb-public"
    : "../alb-private"
  )
}

input = {
  alb_id = dependency.alb.outputs.id
}
```

In the child `terragrunt.hcl`, the `dependency` path for the `alb` depends on whether the VPC is the `mgmt` VPC or not, which is determined by the `dependency.vpc` in the root config. This means that the output from `dependency.vpc` must be available to parse the `dependency.alb` config.

This causes problems when performing a `run-all apply` operation. During a `run-all` operation, Terragrunt first parses all the `dependency` blocks to build a `dependency` tree of the Terragrunt modules to figure out the order of operations. If all the paths are static references, then Terragrunt can determine all the `dependency` paths before any module has been applied. In this case there is no problem even if other config blocks access dependency, as by the time Terragrunt needs to parse those blocks, the upstream dependencies would have been applied during the `run-all apply`.

However, if those `dependency` blocks depend on upstream dependencies, then there is a problem as Terragrunt would not be able to build the dependency tree without the upstream dependencies being applied.

Therefore, to ensure that Terragrunt can build the dependency tree in a `run-all` operation, Terragrunt enforces the following limitation to exposed `include` config:

If the included configuration has any `dependency` blocks, only `locals` and `include` are exposed and available to the child `include` and `dependency` blocks. There are no restrictions for other blocks in the child config (e.g., you can reference `inputs` from the included config in child `inputs`).

Otherwise, if the included config has no `dependency` blocks, there is no restriction on which exposed attributes you can access.

For example, the following alternative configuration is valid even if the alb dependency is still accessing the `inputs` attribute from the included config:

*root terragrunt.hcl*

```hcl
inputs = {
  vpc_name = "mgmt"
}
```

*child terragrunt.hcl*

```hcl
include "root" {
  path   = find_in_parent_folders()
  expose = true
}

dependency "vpc" {
  config_path = "../vpc"
}

dependency "alb" {
  config_path = (
    include.root.inputs.vpc_name == "mgmt"
    ? "../alb-public"
    : "../alb-private"
  )
}

input = {
  vpc_name = dependency.vpc.outputs.name
  alb_id   = dependency.alb.outputs.id
}
```

**What is deep merge?**

When the `merge_strategy` for the `include` block is set to `deep`, Terragrunt will perform a deep merge of the included config. For Terragrunt config, deep merge is defined as follows:

* For simple types, the child overrides the parent.
* For lists, the two attribute lists are combined together in concatenation.
* For maps, the two maps are combined together recursively. That is, if the map keys overlap, then a deep merge is performed on the map value.
* For blocks, if the label is the same, the two blocks are combined together recursively. Otherwise, the blocks are appended like a list. This is similar to maps, with block labels treated as keys.

However, due to internal implementation details, some blocks are not deep mergeable. This will change in the future, but for now, terragrunt performs a shallow merge (that is, block definitions in the child completely override the parent definition). The following blocks have this limitation: - `remote_state` - `generate`

Similarly, the `locals` block is deliberately omitted from the merge operation by design. That is, you will not be able to access parent config `locals` in the child config, and vice versa in a merge. However, you can access the parent locals in child config if you use the `expose` feature.

Finally, `dependency` blocks have special treatment. When doing a `deep` merge, `dependency` blocks from **both** child and parent config are accessible in **both** places. For example, consider the following setup:

*parent config*

```hcl
dependency "vpc" {
  config_path = "../vpc"
}

inputs = {
  vpc_id = dependency.vpc.outputs.vpc_id
  db_id = dependency.mysql.outputs.db_id
}
```

*child config*

```hcl
include "root" {
  path           = find_in_parent_folders()
  merge_strategy = "deep"
}

dependency "mysql" {
  config_path = "../mysql"
}

inputs = {
  security_group_id = dependency.vpc.outputs.security_group_id
}
```

In the example, note how the parent is accessing the outputs of the `mysql` dependency even though it is not defined in the parent. Similarly, the child is accessing the outputs of the `vpc` dependency even though it is not defined in the child.

Full example:

*parent terragrunt.hcl*

```hcl
remote_state {
  backend = "s3"
  config = {
    encrypt = true
    bucket = "__FILL_IN_BUCKET_NAME__"
    key = "${path_relative_to_include()}/terraform.tfstate"
    region = "us-west-2"
  }
}

dependency "vpc" {
  # This will get overridden by child terragrunt.hcl configs
  config_path = ""

  mock_outputs = {
    attribute     = "hello"
    old_attribute = "old val"
    list_attr     = ["hello"]
    map_attr = {
      foo = "bar"
    }
  }
  mock_outputs_allowed_terraform_commands = ["apply", "plan", "destroy", "output"]
}

inputs = {
  attribute     = "hello"
  old_attribute = "old val"
  list_attr     = ["hello"]
  map_attr = {
    foo = "bar"
    test = dependency.vpc.outputs.new_attribute
  }
}
```

*child terragrunt.hcl*

```hcl
include "root" {
  path           = find_in_parent_folders()
  merge_strategy = "deep"
}

remote_state {
  backend = "local"
}

dependency "vpc" {
  config_path = "../vpc"
  mock_outputs = {
    attribute     = "mock"
    new_attribute = "new val"
    list_attr     = ["mock"]
    map_attr = {
      bar = "baz"
    }
  }
}

inputs = {
  attribute     = "mock"
  new_attribute = "new val"
  list_attr     = ["mock"]
  map_attr = {
    bar = "baz"
  }

  dep_out = dependency.vpc.outputs
}
```

*merged terragrunt.hcl*

```hcl
# Child override parent completely due to deep merge limitation
remote_state {
  backend = "local"
}

# mock_outputs are merged together with deep merge
dependency "vpc" {
  config_path = "../vpc"       # Child overrides parent
  mock_outputs = {
    attribute     = "mock"     # Child overrides parent
    old_attribute = "old val"  # From parent
    new_attribute = "new val"  # From child
    list_attr     = [
      "hello",                 # From parent
      "mock",                  # From child
    ]
    map_attr = {
      foo = "bar"              # From parent
      bar = "baz"              # From child
    }
  }

  # From parent
  mock_outputs_allowed_terraform_commands = ["apply", "plan", "destroy", "output"]
}

# inputs are merged together with deep merge
inputs = {
  attribute     = "mock"       # Child overrides parent
  old_attribute = "old val"    # From parent
  new_attribute = "new val"    # From child
  list_attr     = [
    "hello",                 # From parent
    "mock",                  # From child
  ]
  map_attr = {
    foo = "bar"                                   # From parent
    bar = "baz"                                   # From child
    test = dependency.vpc.outputs.new_attribute   # From parent, referencing dependency mock output from child
  }

  dep_out = dependency.vpc.outputs                # From child
}
```
