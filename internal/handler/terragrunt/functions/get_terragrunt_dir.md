get_terragrunt_dir()
---
`get_terragrunt_dir()` returns the directory where the Terragrunt configuration file (by default `terragrunt.hcl`) lives. This is useful when you need to use relative paths with [remote Terraform configurations]({{site.baseurl}}/docs/features/keep-your-terraform-code-dry/#remote-terraform-configurations) and you want those paths relative to your Terragrunt configuration file and not relative to the temporary directory where Terragrunt downloads the code.

For example, imagine you have the following file structure:

    /terraform-code
    ├── common.tfvars
    ├── frontend-app
    │   └── terragrunt.hcl

Inside of `/terraform-code/frontend-app/terragrunt.hcl` you might try to write code that looks like this:

``` hcl
terraform {
  source = "git::git@github.com:foo/modules.git//frontend-app?ref=v0.0.3"

  extra_arguments "custom_vars" {
    commands = [
      "apply",
      "plan",
      "import",
      "push",
      "refresh"
    ]

    arguments = [
      "-var-file=../common.tfvars" # Note: This relative path will NOT work correctly!
    ]
  }
}
```

Note how the `source` parameter is set, so Terragrunt will download the `frontend-app` code from the `modules` repo into a temporary folder and run `terraform` in that temporary folder. Note also that there is an `extra_arguments` block that is trying to allow the `frontend-app` to read some shared variables from a `common.tfvars` file. Unfortunately, the relative path (`../common.tfvars`) won’t work, as it will be relative to the temporary folder\! Moreover, you can’t use an absolute path, or the code won’t work on any of your teammates' computers.

To make the relative path work, you need to use `get_terragrunt_dir()` to combine the path with the folder where the `terragrunt.hcl` file lives:

``` hcl
terraform {
  source = "git::git@github.com:foo/modules.git//frontend-app?ref=v0.0.3"

  extra_arguments "custom_vars" {
    commands = [
      "apply",
      "plan",
      "import",
      "push",
      "refresh"
    ]

    # With the get_terragrunt_dir() function, you can use relative paths!
    arguments = [
      "-var-file=${get_terragrunt_dir()}/../common.tfvars"
    ]
  }
}
```
