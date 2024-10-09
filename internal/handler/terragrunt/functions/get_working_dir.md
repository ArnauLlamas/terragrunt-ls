get_working_dir()
---
`get_working_dir()` returns the absolute path where Terragrunt runs Terraform commands. This is useful when you need to manage substitutions of vars inside a *.tfvars file located right inside terragrunt's tmp dir.
