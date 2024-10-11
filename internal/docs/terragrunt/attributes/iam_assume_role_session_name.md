iam_assume_role_session_name
---
The `iam_assume_role_session_name` attribute can be used to specify the STS session name, for the IAM role that Terragrunt should assume prior to invoking Terraform.

The precedence is as follows: `--terragrunt-iam-assume-role-session-name` command line option → `TERRAGRUNT_IAM_ASSUME_ROLE_SESSION_NAME` env variable →
`iam_assume_role_session_name` attribute of the `terragrunt.hcl` file in the module directory → `iam_assume_role_session_name` attribute of the included
`terragrunt.hcl`.
