get_terragrunt_source_cli_flag()
---
`get_terragrunt_source_cli_flag()` returns the value passed in via the CLI `--terragrunt-source` or an environment variable `TERRAGRUNT_SOURCE`. Note that this will return an empty string when either of those values are not provided.

This is useful for constructing before and after hooks, or TF flags that only apply to local development (e.g., setting up debug flags, or adjusting the `iam_role` parameter).

Some example use cases are:

- Setting debug logging when doing local development.
- Adjusting the kubernetes provider configuration so that it targets minikube instead of real clusters.
- Providing special mocks pulled in from the local dev source (e.g., something like `mock_outputs = jsondecode(file("${get_terragrunt_source_cli_arg()}/dependency_mocks/vpc.json"))`).
