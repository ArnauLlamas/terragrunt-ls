get_terraform_commands_that_need_vars()
---
`get_terraform_commands_that_need_vars()` returns the list of terraform commands that accept `-var` and `-var-file` parameters. This function is used when defining [extra\_arguments]({{site.baseurl}}/docs/features/keep-your-cli-flags-dry/#multiple-extra_arguments-blocks).

``` hcl
terraform {
  extra_arguments "common_var" {
    commands  = get_terraform_commands_that_need_vars()
    arguments = ["-var-file=${get_aws_account_id()}.tfvars"]
  }
}
```
