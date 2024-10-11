get_platform()
---
`get_platform()` returns the current Operating System. Example:

``` hcl
inputs = {
  platform = get_platform()
}
```

This function can also be used in a comparison to evaluate what to do based on the current operating system. Example:
``` hcl
output "platform" {
  value = var.platform == "darwin" ? "(value for MacOS)" : "(value for other OS's)"
}
```

Some of the returned values can be:
```
darwin
freebsd
linux
windows
```
