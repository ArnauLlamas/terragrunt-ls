get_env($name)
---
`get_env(NAME)` return the value of variable named `NAME` or throws exceptions if that variable is not set. Example:

``` hcl
remote_state {
  backend = "s3"
  config = {
    bucket = get_env("BUCKET")
  }
}
```

`get_env(NAME, DEFAULT)` returns the value of the environment variable named `NAME` or `DEFAULT` if that environment variable is not set. Example:

``` hcl
remote_state {
  backend = "s3"
  config = {
    bucket = get_env("BUCKET", "my-terraform-bucket")
  }
}
```

Note that [Terraform will read environment variables](https://www.terraform.io/docs/configuration/environment-variables.html#tf_var_name) that start with the prefix `TF_VAR_`, so one way to share a variable named `foo` between Terraform and Terragrunt is to set its value as the environment variable `TF_VAR_foo` and to read that value in using this `get_env()` built-in function.
