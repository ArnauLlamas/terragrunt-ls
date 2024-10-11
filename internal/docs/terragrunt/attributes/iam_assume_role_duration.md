iam_assume_role_duration
---
The `iam_assume_role_duration` attribute can be used to specify the STS session duration, in seconds, for the IAM role that Terragrunt should assume prior to invoking Terraform.

The precedence is as follows: `--terragrunt-iam-assume-role-duration` command line option → `TERRAGRUNT_IAM_ASSUME_ROLE_DURATION` env variable →
`iam_assume_role_duration` attribute of the `terragrunt.hcl` file in the module directory → `iam_assume_role_duration` attribute of the included
`terragrunt.hcl`.

Example:

```hcl
iam_assume_role_duration = 14400
```
