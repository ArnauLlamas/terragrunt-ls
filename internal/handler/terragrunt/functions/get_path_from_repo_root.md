get_path_from_repo_root()
---
`get_path_from_repo_root()` returns the path from the root of the Git repository to the current directory:

```hcl
remote_state {
  backend = "s3"

  config = {
    bucket         = "terraform"
    dynamodb_table = "terraform"
    encrypt        = true
    key            = "${get_path_from_repo_root()}/terraform.tfstate"
    session_name   = "terraform"
    region         = "us-east-1"
  }
}
```

This function will error if the file is not located in a Git repository.
