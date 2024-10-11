get_repo_root()
---
Returns the absolute path to the root of the Git repository:

This function will error if the file is not located in a Git repository.

Example:

```hcl
inputs {
    very_important_config = "${get_repo_root()}/config/strawberries.conf"
}
```
