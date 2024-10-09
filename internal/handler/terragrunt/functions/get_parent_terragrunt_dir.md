get_parent_terragrunt_dir($name)
---
`get_parent_terragrunt_dir()` returns the absolute directory where the Terragrunt parent configuration file (by default `terragrunt.hcl`) lives. This is useful when you need to use relative paths with [remote Terraform configurations]({{site.baseurl}}/docs/features/keep-your-terraform-code-dry/#remote-terraform-configurations) and you want those paths relative to your parent Terragrunt configuration file and not relative to the temporary directory where Terragrunt downloads the code.

This function is very similar to [get\_terragrunt\_dir()](#get_terragrunt_dir) except it returns the root instead of the leaf of your terragrunt configuration folder.

    /terraform-code
    ├── terragrunt.hcl
    ├── common.tfvars
    ├── app1
    │   └── terragrunt.hcl
    ├── tests
    │   ├── app2
    │   |   └── terragrunt.hcl
    │   └── app3
    │       └── terragrunt.hcl

``` hcl
terraform {
  extra_arguments "common_vars" {
    commands = [
      "apply",
      "plan",
      "import",
      "push",
      "refresh"
    ]

    arguments = [
      "-var-file=${get_parent_terragrunt_dir()}/common.tfvars"
    ]
  }
}
```

The common.tfvars located in the terraform root folder will be included by all applications, whatever their relative location to the root.

If you have `include` blocks, this function requires a `name` parameter when used in the child config to specify which
`include` block to base the parent dir on.

Example:

```hcl
include "root" {
  path = find_in_parent_folders()
}
include "region" {
  path = find_in_parent_folders("region.hcl")
}

terraform {
  source = "${get_parent_terragrunt_dir("root")}/modules/vpc"
}
```
