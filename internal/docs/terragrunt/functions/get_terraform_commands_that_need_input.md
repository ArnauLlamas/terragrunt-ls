get_terraform_commands_that_need_input()
---
`get_terraform_commands_that_need_input()` returns the list of terraform commands that accept the `-input=(true or false)` parameter. This function is used when defining [extra\_arguments]({{site.baseurl}}/docs/features/keep-your-cli-flags-dry/#multiple-extra_arguments-blocks).

``` hcl
terraform {
  # Force Terraform to not ask for input value if some variables are undefined.
  extra_arguments "disable_input" {
    commands  = get_terraform_commands_that_need_input()
    arguments = ["-input=false"]
  }
}
```
