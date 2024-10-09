run_cmd($command, )
---
`run_cmd(command, arg1, arg2...)` runs a shell command and returns the stdout as the result of the interpolation. The command is executed at the same folder as the `terragrunt.hcl` file. This is useful whenever you want to dynamically fill in arbitrary information in your Terragrunt configuration.

As an example, you could write a script that determines the bucket and DynamoDB table name based on the AWS account, instead of hardcoding the name of every account:

``` hcl
remote_state {
  backend = "s3"
  config = {
    bucket         = run_cmd("./get_names.sh", "bucket")
    dynamodb_table = run_cmd("./get_names.sh", "dynamodb")
  }
}
```

If the command you are running has the potential to output sensitive values, you may wish to redact the output from appearing in the terminal. To do so, use the special `--terragrunt-quiet` argument which must be passed as one of the first arguments to `run_cmd()`:

``` hcl
super_secret_value = run_cmd("--terragrunt-quiet", "./decrypt_secret.sh", "foo")
```

**Note:** This will prevent terragrunt from displaying the output from the command in its output. However, the value could still be displayed in the Terraform output if Terraform does not treat it as a [sensitive value](https://www.terraform.io/docs/configuration/outputs.html#sensitive-suppressing-values-in-cli-output).

Invocations of `run_cmd` are cached based on directory and executed command, so cached values are re-used later, rather than executed multiple times. Here's an example:
```hcl
locals {
  uuid = run_cmd("echo", "uuid1",  uuid())
  uuid2 = run_cmd("echo", "uuid2", uuid())
  uuid3 = run_cmd("echo", "uuid3", uuid())
  potato = run_cmd("echo", "potato")
  potato2 = run_cmd("echo", "potato")
  carrot = run_cmd("echo", "carrot")
}
inputs = {
  potato3 = run_cmd("echo", "potato")
  uuid3 = run_cmd("echo", "uuid3", uuid())
  uuid4 = run_cmd("echo", "uuid4", uuid())
  carrot2 = run_cmd("echo", "carrot")
}
```

Output:
```
$ terragrunt init
uuid1 b48379e1-924d-2403-8789-c72d50be964c
uuid1 9f3a8398-b11f-5314-7783-dad176ee487d
uuid1 649ac501-e5db-c935-1499-c59fb7a75625
uuid2 2d65972b-3fa9-181f-64fe-dcd574d944d0
uuid3 e345de60-9cfa-0455-79b7-af0d053a15a5
potato
uuid3 7f90a4ed-96e3-1dd8-5fee-91b8c8e07650
uuid2 8638fe79-c589-bebd-2a2a-3e6b96f7fc34
uuid3 310d0447-f0a6-3f67-efda-e6b1521fa1fb
uuid4 f8e80cc6-1892-8db7-bd63-6089fef00c01
uuid2 289ff371-8021-54c6-2254-72de9d11392a
uuid3 baa19863-1d99-e0ef-11f2-ede830d1c58a
carrot
```
**Notes:**
  * Output contains only once `carrot` and `potato`, because other invocations got cached, caching works for all sections
  * Output contains multiple times `uuid1` and `uuid2` because during HCL evaluation each `run_cmd` in `locals` is evaluated multiple times and random argument generated from `uuid()` save cached value under different key each time
  * Output contains multiple times `uuid3`, +1 more output comparing to `uuid1` and `uuid2` - because `uuid3` is declared in locals and inputs which add one more evaluation
  * Output contains only once `uuid4` since it is declared only once in `inputs`, `inputs` is not evaluated twice

You can modify this caching behavior to ignore the existing directory if you know the command you are running is not dependent on the current directory path.  To do so, use the special `--terragrunt-global-cache` argument which must be passed as one of the first arguments to `run_cmd()` (and can be combined with `--terragrunt-quiet` in any order):

``` hcl
value = run_cmd("--terragrunt-global-cache", "--terragrunt-quiet", "/usr/local/bin/get-account-map")
```
