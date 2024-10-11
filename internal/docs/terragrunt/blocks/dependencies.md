dependencies
---
The `dependencies` block is used to enumerate all the Terragrunt modules that need to be applied in order for this module to be able to apply. Note that this is purely for ordering the operations when using `run-all` commands of Terraform. This does not expose or pull in the outputs like `dependency` blocks.

The `dependencies` block supports the following arguments:
* `paths` (attribute): A list of paths to modules that should be marked as a dependency.

Example:

```hcl
# When applying this terragrunt config in an `run-all` command, make sure the modules at "../vpc" and "../rds" are
# handled first.
dependencies {
  paths = ["../vpc", "../rds"]
}
```
