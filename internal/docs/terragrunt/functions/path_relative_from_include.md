path_relative_from_include($name)
---
`path_relative_from_include()` returns the relative path between the `path` specified in its `include` block and the current `terragrunt.hcl` file (it is the counterpart of `path_relative_to_include()`). For example, consider the following folder structure:

    ├── sources
    |  ├── mysql
    |  |  └── \*.tf
    |  └── secrets
    |     └── mysql
    |         └── \*.tf
    └── terragrunt
      └── common.tfvars
      ├── mysql
      |  └── terragrunt.hcl
      ├── secrets
      |  └── mysql
      |     └── terragrunt.hcl
      └── terragrunt.hcl

Imagine `terragrunt/mysql/terragrunt.hcl` and `terragrunt/secrets/mysql/terragrunt.hcl` include all settings from the root `terragrunt.hcl` file:

``` hcl
include "root" {
  path = find_in_parent_folders()
}
```

The root `terragrunt.hcl` can use the `path_relative_from_include()` in combination with `path_relative_to_include()` in its `source` configuration to retrieve the relative terraform source code from the terragrunt configuration file:

``` hcl
terraform {
  source = "${path_relative_from_include()}/../sources//${path_relative_to_include()}"
}
```

The resulting `source` will be `../../sources//mysql` for `mysql` module and `../../../sources//secrets/mysql` for `secrets/mysql` module.

Another use case would be to add extra argument to include the `common.tfvars` file for all subdirectories:

``` hcl
  terraform {
    extra_arguments "common_var" {
      commands = [
        "apply",
        "plan",
        "import",
        "push",
        "refresh"
      ]

      arguments = [
        "-var-file=${get_terragrunt_dir()}/${path_relative_from_include()}/common.tfvars",
      ]
    }
  }
```

This allows proper retrieval of the `common.tfvars` from whatever the level of subdirectories we have.

If you have `include` blocks, this function requires a `name` parameter when used in the child config to specify which
`include` block to base the relative path on.

Example:

```hcl
include "root" {
  path = find_in_parent_folders()
}
include "region" {
  path = find_in_parent_folders("region.hcl")
}

terraform {
  source = "../modules/${path_relative_from_include("root")}"
}
```
