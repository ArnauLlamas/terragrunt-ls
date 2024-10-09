read_tfvars_file($file)
---
`read_tfvars_file(file_path)` reads a `.tfvars` or `.tfvars.json` file and returns a map of the variables defined in it.

This is useful for reading variables from a `.tfvars` file and merging them into the inputs or to use them in a `locals` block:

```hcl

locals {
  inputs_from_tfvars = jsondecode(read_tfvars_file("common.tfvars"))
}

inputs = merge(
  local.inputs_from_tfvars,
  {
    # additional inputs
  }
)
```

Another example:

```hcl

locals {
  backend = jsondecode(read_tfvars_file("backend.tfvars"))
}

remote_state {
  backend = "s3"
  config = {
    bucket         = "${get_env("TG_BUCKET_PREFIX", "tf-bucket")}-${get_aws_account_id()}"
    key            = "${path_relative_to_include()}/terraform-${local.aws_region}.tfstate"
    region         = local.backend.region
  }
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite_terragrunt"
  }
}
```
