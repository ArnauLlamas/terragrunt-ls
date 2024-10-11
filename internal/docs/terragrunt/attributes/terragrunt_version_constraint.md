terragrunt_version_constraint
---
The terragrunt `terragrunt_version_constraint` string can be used to specify which versions of the Terragrunt CLI can be used with your configuration. If the running version of Terragrunt doesn't match the constraints specified, Terragrunt will produce an error and exit without taking any further actions.

Example:

```hcl
terragrunt_version_constraint = ">= 0.23"
```
