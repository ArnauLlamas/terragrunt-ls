download_dir
---
The terragrunt `download_dir` string option can be used to override the default download directory.

The precedence is as follows: `--terragrunt-download-dir` command line option → `TERRAGRUNT_DOWNLOAD` env variable →
`download_dir` attribute of the `terragrunt.hcl` file in the module directory → `download_dir` attribute of the included
`terragrunt.hcl`.

It supports all terragrunt functions, i.e. `path_relative_from_include()`.
