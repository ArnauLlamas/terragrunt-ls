iam_role
---
The `iam_role` attribute can be used to specify an IAM role that Terragrunt should assume prior to invoking Terraform.

The precedence is as follows: `--terragrunt-iam-role` command line option → `TERRAGRUNT_IAM_ROLE` env variable →
`iam_role` attribute of the `terragrunt.hcl` file in the module directory → `iam_role` attribute of the included
`terragrunt.hcl`.

Example:

```hcl
iam_role = "arn:aws:iam::ACCOUNT_ID:role/ROLE_NAME"
```
**Notes:**
  * Value of `iam_role` can reference local variables
  * Definitions of `iam_role` included from other HCL files through `include`
