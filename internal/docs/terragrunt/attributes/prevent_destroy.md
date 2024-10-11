prevent_destroy
---
Terragrunt `prevent_destroy` boolean flag allows you to protect selected Terraform module. It will prevent `destroy` or
`destroy-all` command to actually destroy resources of the protected module. This is useful for modules you want to
carefully protect, such as a database, or a module that provides auth.

Example:

```hcl
terraform {
  source = "git::git@github.com:foo/modules.git//app?ref=v0.0.3"
}

prevent_destroy = true
```
