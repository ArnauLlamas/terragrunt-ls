get_terraform_commands_that_need_parallelism()
---
`get_terraform_commands_that_need_parallelism()` returns the list of terraform commands that accept the `-parallelism` parameter. This function is used when defining [extra\_arguments]({{site.baseurl}}/docs/features/keep-your-cli-flags-dry/#multiple-extra_arguments-blocks).

``` hcl
terraform {
  # Force Terraform to run with reduced parallelism
  extra_arguments "parallelism" {
    commands  = get_terraform_commands_that_need_parallelism()
    arguments = ["-parallelism=5"]
  }
}
```
