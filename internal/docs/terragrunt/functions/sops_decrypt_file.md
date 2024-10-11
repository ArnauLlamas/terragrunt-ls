sops_decrypt_file($file)
---
`sops_decrypt_file(file_path)` decrypts a yaml, json, ini, env or "raw text" file encrypted with `sops`.

[sops](https://github.com/mozilla/sops) is an editor of encrypted files that supports YAML, JSON, ENV, INI and
BINARY formats and encrypts with AWS KMS, GCP KMS, Azure Key Vault, Hashicorp Vault and PGP.

This allows static secrets to be stored encrypted within your Terragrunt repository.

For example, suppose you have some static secrets required to bootstrap your
infrastructure in `secrets.yaml`, you can decrypt and merge them into the inputs
by using `sops_decrypt_file`:

```hcl
locals {
  secret_vars = yamldecode(sops_decrypt_file(find_in_parent_folders("secrets.yaml")))
}

inputs = merge(
  local.secret_vars,
  {
    # additional inputs
  }
)
```

If you absolutely need to fallback to a default value you can make use of the Terraform `try` function:

```hcl
locals {
  secret_vars = try(jsondecode(sops_decrypt_file(find_in_parent_folders("no-secrets-here.json"))), {})
}

inputs = merge(
  local.secret_vars, # This will be {}
  {
    # additional inputs
  }
)
```
