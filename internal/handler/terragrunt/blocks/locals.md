locals
---
The `locals` block is used to define aliases for Terragrunt expressions that can be referenced within the configuration. You can learn more about `locals` in the feature overview.

The `locals` block does not have a defined set of arguments that are supported. Instead, all the arguments passed into `locals` are available under the reference `local.ARG_NAME` throughout the Terragrunt configuration.

Example:

```hcl
# Make the AWS region a reusable variable within the configuration
locals {
  aws_region = "us-east-1"
}

inputs = {
  region = local.aws_region
  name   = "${local.aws_region}-bucket"
}
```
