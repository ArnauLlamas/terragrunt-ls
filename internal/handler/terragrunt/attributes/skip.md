skip
---
The terragrunt `skip` boolean flag can be used to protect modules you don’t want any changes to or just to skip modules
that don’t define any infrastructure by themselves. When set to true, all terragrunt commands will skip the selected
module.

Consider the following file structure:

    root
    ├── terragrunt.hcl
    ├── prod
    │   └── terragrunt.hcl
    ├── dev
    │   └── terragrunt.hcl
    └── qa
        └── terragrunt.hcl

In some cases, the root level `terragrunt.hcl` file is solely used to DRY up your Terraform configuration by being
included in the other `terragrunt.hcl` files. In this case, you do not want the `run-all` commands to process the root
level `terragrunt.hcl` since it does not define any infrastructure by itself. To make the `run-all` commands skip the
root level `terragrunt.hcl` file, you can set `skip = true`:

``` hcl
skip = true
```

The `skip` flag must be set explicitly in terragrunt modules that should be skipped. If you set `skip = true` in a
`terragrunt.hcl` file that is included by another `terragrunt.hcl` file, only the `terragrunt.hcl` file that explicitly
set `skip = true` will be skipped.
