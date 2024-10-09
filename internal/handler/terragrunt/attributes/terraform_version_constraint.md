terraform_version_constraint
---
The terragrunt `terraform_version_constraint` string overrides the default minimum supported version of terraform.
Terragrunt only officially supports the latest version of terraform, however in some cases an old terraform is needed.

Example:

```hcl
terraform_version_constraint = ">= 0.11"
```
