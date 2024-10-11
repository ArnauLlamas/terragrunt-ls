get_aws_caller_identity_user_id()
---
`get_aws_caller_identity_user_id()` returns the UserId of the AWS identity associated with the current set of credentials. Example:

``` hcl
inputs = {
  caller_user_id = get_aws_caller_identity_user_id()
}
```

This allows uniqueness of the storage bucket per AWS account (since bucket name must be globally unique).

It is also possible to configure variables specifically based on the account used:

``` hcl
terraform {
  extra_arguments "common_var" {
    commands = get_terraform_commands_that_need_vars()
    arguments = ["-var-file=${get_aws_account_id()}.tfvars"]
  }
}
```
**Note:** value returned by `get_aws_caller_identity_user_id()` can change during parsing of HCL code, for example after evaluation of `iam_role` attribute.
