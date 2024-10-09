get_default_retryable_errors()
---
`get_default_retryable_errors()` returns default retryabled errors. Example:

``` hcl
retryable_errors = concat(get_default_retryable_errors(), ["my custom error"])
```
