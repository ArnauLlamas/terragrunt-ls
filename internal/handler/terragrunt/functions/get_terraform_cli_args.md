get_terraform_cli_args()
---
`get_terraform_cli_args()` returns cli args for the current terraform command in execution. Example:

``` hcl
inputs = {
  current_cli_args = get_terraform_cli_args()
}
```
