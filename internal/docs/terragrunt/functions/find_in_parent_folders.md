find_in_parent_folders($file)
---
`find_in_parent_folders()` searches up the directory tree from the current `terragrunt.hcl` file and returns the absolute path to the first `terragrunt.hcl` in a parent folder or exit with an error if no such file is found. This is primarily useful in an `include` block to automatically find the path to a parent `terragrunt.hcl` file:

``` hcl
include "root" {
  path = find_in_parent_folders()
}
```

The function takes an optional `name` parameter that allows you to specify a different filename to search for:

``` hcl
include "root" {
  path = find_in_parent_folders("some-other-file-name.hcl")
}
```

You can also pass an optional second `fallback` parameter which causes the function to return the fallback value (instead of exiting with an error) if the file in the `name` parameter cannot be found:

``` hcl
include "root" {
  path = find_in_parent_folders("some-other-file-name.hcl", "fallback.hcl")
}
```

Note that this function searches relative to the child `terragrunt.hcl` file when called from a parent config. For
example, if you had the following folder structure:

    ├── terragrunt.hcl
    └── prod
        ├── env.hcl
        └── mysql
            └── terragrunt.hcl

And the root `terragrunt.hcl` contained the following:

    locals {
      env_vars = read_terragrunt_config(find_in_parent_folders("env.hcl"))
    }

The `find_in_parent_folders` will search from the __child `terragrunt.hcl`__ (`prod/mysql/terragrunt.hcl`) config,
finding the `env.hcl` file in the `prod` directory.
