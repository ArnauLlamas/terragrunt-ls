path_relative_to_include($name)
---
`path_relative_to_include()` returns the relative path between the current `terragrunt.hcl` file and the `path` specified in its `include` block. For example, consider the following folder structure:

    ├── terragrunt.hcl
    └── prod
        └── mysql
            └── terragrunt.hcl
    └── stage
        └── mysql
            └── terragrunt.hcl

Imagine `prod/mysql/terragrunt.hcl` and `stage/mysql/terragrunt.hcl` include all settings from the root `terragrunt.hcl` file:

``` hcl
include "root" {
  path = find_in_parent_folders()
}
```

The root `terragrunt.hcl` can use the `path_relative_to_include()` in its `remote_state` configuration to ensure each child stores its remote state at a different `key`:

``` hcl
remote_state {
  backend = "s3"
  config = {
    bucket = "my-terraform-bucket"
    region = "us-east-1"
    key    = "${path_relative_to_include()}/terraform.tfstate"
  }
}
```

The resulting `key` will be `prod/mysql/terraform.tfstate` for the prod `mysql` module and `stage/mysql/terraform.tfstate` for the stage `mysql` module.

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
  source = "../modules/${path_relative_to_include("root")}"
}
```
