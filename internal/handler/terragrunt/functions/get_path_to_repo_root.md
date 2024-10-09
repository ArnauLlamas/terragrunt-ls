get_path_to_repo_root()
---
`get_path_to_repo_root()` returns the relative path to the root of the Git repository:

```hcl
terraform {
  source = "${get_path_to_repo_root()}//modules/example"
}
```

This function will error if the file is not located in a Git repository.
