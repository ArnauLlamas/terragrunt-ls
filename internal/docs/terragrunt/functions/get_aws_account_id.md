get_aws_account_id()
---
`get_aws_account_id()` returns the AWS account id associated with the current set of credentials. Example:

``` hcl
remote_state {
  backend = "s3"
  config = {
    bucket = "mycompany-${get_aws_account_id()}"
  }
}
```
**Note:** value returned by `get_aws_account_id()` can change during parsing of HCL code, for example after evaluation of `iam_role` attribute.
