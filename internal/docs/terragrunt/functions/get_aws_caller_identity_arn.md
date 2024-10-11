get_aws_caller_identity_arn()
---
`get_aws_caller_identity_arn()` returns the ARN of the AWS identity associated with the current set of credentials. Example:

``` hcl
inputs = {
  caller_arn = get_aws_caller_identity_arn()
}
```
**Note:** value returned by `get_aws_caller_identity_arn()` can change during parsing of HCL code, for example after evaluation of `iam_role` attribute.
