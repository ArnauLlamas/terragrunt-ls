get_original_terragrunt_dir()
---
`get_original_terragrunt_dir()` returns the directory where the original Terragrunt configuration file (by default
`terragrunt.hcl`) lives. This is primarily useful when one Terragrunt config is being read from another: e.g., if
`/terraform-code/terragrunt.hcl` calls `read_terragrunt_config("/foo/bar.hcl")`, and within `bar.hcl`, you call
`get_original_terragrunt_dir()`, you'll get back `/terraform-code`.
