read_terragrunt_config($file, $fallback)
---
`read_terragrunt_config(config_path, [default_val])` parses the terragrunt config at the given path and serializes the
result into a map that can be used to reference the values of the parsed config. This function will expose all blocks
and attributes of a terragrunt config.

For example, suppose you had a config file called `common.hcl` that contains common input variables:

```hcl
inputs = {
  stack_name = "staging"
  account_id = "1234567890"
}
```

You can read these inputs in another config by using `read_terragrunt_config`, and merge them into the inputs:

```hcl
locals {
  common_vars = read_terragrunt_config(find_in_parent_folders("common.hcl"))
}

inputs = merge(
  local.common_vars.inputs,
  {
    # additional inputs
  }
)
```

This function also takes in an optional second parameter which will be returned if the file does not exist:

```hcl
locals {
  common_vars = read_terragrunt_config(find_in_parent_folders("i-dont-exist.hcl", "i-dont-exist.hcl"), {inputs = {}})
}

inputs = merge(
  local.common_vars.inputs, # This will be {}
  {
    # additional inputs
  }
)
```

Note that this function will also render `dependency` blocks. That is, the parsed config will make the outputs of the
`dependency` blocks available. For example, suppose you had the following config in a file called `common_deps.hcl`:

```hcl
dependency "vpc" {
  config_path = "${get_terragrunt_dir()}/../vpc"
}
```

You can access the outputs of the vpc dependency through the parsed outputs of `read_terragrunt_config`:

```hcl
locals {
  common_deps = read_terragrunt_config(find_in_parent_folders("common_deps.hcl"))
}

inputs = {
  vpc_id = local.common_deps.dependency.vpc.outputs.vpc_id
}
```
