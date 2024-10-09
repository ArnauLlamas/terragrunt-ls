get_terraform_command()
---
`get_terraform_command()` returns the current terraform command in execution. Example:

``` hcl
inputs = {
  current_command = get_terraform_command()
}
```
