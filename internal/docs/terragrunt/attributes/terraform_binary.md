terraform_binary
---
The terragrunt `terraform_binary` string option can be used to override the default terraform binary path (which is
`terraform`).

The precedence is as follows: `--terragrunt-tfpath` command line option → `TERRAGRUNT_TFPATH` env variable →
`terragrunt.hcl` in the module directory → included `terragrunt.hcl`
