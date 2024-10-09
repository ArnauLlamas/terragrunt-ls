terraform
---
The `terraform` block is used to configure how Terragrunt will interact with Terraform. This includes specifying where
to find the Terraform configuration files, any extra arguments to pass to the `terraform` CLI, and any hooks to run
before or after calling Terraform.

The `terraform` block supports the following arguments:

- `source` (attribute): Specifies where to find Terraform configuration files. This parameter supports the exact same syntax as the
  [module source](https://www.terraform.io/docs/modules/sources.html) parameter for Terraform `module` blocks **except
  for the Terraform registry** (see below note), including local file paths, Git URLs, and Git URLS with `ref`
  parameters. Terragrunt will download all the code in the repo (i.e. the part before the double-slash `//`) so that
  relative paths work correctly between modules in that repo.
    - The `source` parameter can be configured to pull Terraform modules from any Terraform module registry using
      the `tfr` protocol. The `tfr` protocol expects URLs to be provided in the format
      `tfr://REGISTRY_HOST/MODULE_SOURCE?version=VERSION`. For example, to pull the `terraform-aws-modules/vpc/aws`
      module from the public Terraform registry, you can use the following as the source parameter:
      `tfr://registry.terraform.io/terraform-aws-modules/vpc/aws?version=3.3.0`.
    - If you wish to access a private module registry (e.g., [Terraform Cloud/Enterprise](https://www.terraform.io/docs/cloud/registry/index.html)),
      you can provide the authentication to Terragrunt as an environment variable with the key `TG_TF_REGISTRY_TOKEN`.
      This token can be any registry API token.
    - The `tfr` protocol supports a shorthand notation where the `REGISTRY_HOST` can be omitted to default to the public
      registry (`registry.terraform.io`) if you use `tfr:///` (note the three `/`). For example, the following will
      fetch the `terraform-aws-modules/vpc/aws` module from the public registry:
      `tfr:///terraform-aws-modules/vpc/aws?version=3.3.0`.
    - You can also use submodules from the registry using `//`. For example, to use the `iam-policy` submodule from the
      registry module
      [terraform-aws-modules/iam](https://registry.terraform.io/modules/terraform-aws-modules/iam/aws/latest), you can
      use the following: `tfr:///terraform-aws-modules/iam/aws//modules/iam-policy?version=4.3.0`.
    - Refer to [A note about using modules from the
      registry]({{site.baseurl}}/docs/getting-started/quick-start#a-note-about-using-modules-from-the-registry) for more
      information about using modules from the Terraform Registry with Terragrunt.

- `include_in_copy` (attribute): A list of glob patterns (e.g., `["*.txt"]`) that should always be copied into the
  Terraform working directory. When you use the `source` param in your Terragrunt config and run `terragrunt <command>`,
  Terragrunt will download the code specified at source into a scratch folder (`.terragrunt-cache`, by default), copy
  the code in your current working directory into the same scratch folder, and then run `terraform <command>` in that
  scratch folder. By default, Terragrunt excludes hidden files and folders during the copy step. This feature allows you
  to specify glob patterns of files that should always be copied from the Terragrunt working directory. Additional
  notes:
    - The path should be specified relative to the source directory.
    - This list is also used when using a local file source (e.g., `source = "../modules/vpc"`). For example, if your
      terraform module source contains a hidden file that you want to copy over (e.g., a `.python-version` file), you
      can specify that in this list to ensure it gets copied over to the scratch copy
      (e.g., `include_in_copy = [".python-version"]`).

- `extra_arguments` (block): Nested blocks used to specify extra CLI arguments to pass to the `terraform` CLI. Learn more
  about its usage in the [Keep your CLI flags DRY]({{site.baseurl}}/docs/features/keep-your-cli-flags-dry/) use case overview. Supports
  the following arguments:
    - `arguments` (required) : A list of CLI arguments to pass to `terraform`.
    - `commands` (required) : A list of `terraform` sub commands that the arguments will be passed to.
    - `env_vars` (optional) : A map of key value pairs to set as environment variables when calling `terraform`.
    - `required_var_files` (optional): A list of file paths to terraform vars files (`.tfvars`) that will be passed in to
      `terraform` as `-var-file=<your file>`.
    - `optional_var_files` (optional): A list of file paths to terraform vars files (`.tfvars`) that will be passed in to
      `terraform` like `required_var_files`, only any files that do not exist are ignored.

- `before_hook` (block): Nested blocks used to specify command hooks that should be run before `terraform` is called.
  Hooks run from the directory with the terraform module, except for hooks related to `terragrunt-read-config` and
  `init-from-module`. These hooks run in the terragrunt configuration directory (the directory where `terragrunt.hcl`
  lives).
  Supports the following arguments:
    - `commands` (required) : A list of `terraform` sub commands for which the hook should run before.
    - `execute` (required) : A list of command and arguments that should be run as the hook. For example, if `execute` is set as
      `["echo", "Foo"]`, the command `echo Foo` will be run.
    - `working_dir` (optional) : The path to set as the working directory of the hook. Terragrunt will switch directory
      to this path prior to running the hook command. Defaults to the terragrunt configuration directory for
      `terragrunt-read-config` and `init-from-module` hooks, and the terraform module directory for other command hooks.
    - `run_on_error` (optional) : If set to true, this hook will run even if a previous hook hit an error, or in the
      case of "after" hooks, if the Terraform command hit an error. Default is false.
    - `suppress_stdout` (optional) : If set to true, the stdout output of the executed commands will be suppressed. This can be useful when there are scripts relying on terraform's output and any other output would break their parsing.

- `after_hook` (block): Nested blocks used to specify command hooks that should be run after `terraform` is called.
  Hooks run from the terragrunt configuration directory (the directory where `terragrunt.hcl` lives). Supports the same
  arguments as `before_hook`.
- `error_hook` (block): Nested blocks used to specify command hooks that run when an error is thrown. The
error must match one of the expressions listed in the `on_errors` attribute. Error hooks are executed after the before/after hooks.

In addition to supporting before and after hooks for all terraform commands, the following specialized hooks are also
supported:

- `terragrunt-read-config` (after hook only): `terragrunt-read-config` is a special hook command that you can use with
  the `after_hook` subblock to run an action immediately after terragrunt finishes loading the config. This hook will
  run on every invocation of terragrunt. Note that you can only use this hook with `after_hooks`. Any `before_hooks`
  with the command `terragrunt-read-config` will be ignored. The working directory for hooks associated with this
  command will be the terragrunt config directory.

- `init-from-module` and `init`: Terragrunt has two stages of initialization: one is to download [remote
  configurations](https://terragrunt.gruntwork.io/use-cases/keep-your-terraform-code-dry/) using `go-getter`; the other
  is [Auto-Init](https://terragrunt.gruntwork.io/docs/features/auto-init/), which configures the backend and downloads
  provider plugins and modules. If you wish to run a hook when Terragrunt is using `go-getter` to download remote
  configurations, use `init-from-module` for the command. If you wish to execute a hook when Terragrunt is using
  terraform `init` for Auto-Init, use `init` for the command. For example, an `after_hook` for the command
  `init-from-module` will run after terragrunt clones the module, while an `after_hook` for the command `init` will run
  after terragrunt runs `terraform init` on the cloned module.
    - Hooks for both `init-from-module` and `init` only run if the requisite stage needs to run. That is, if terragrunt
      detects that the module is already cloned in the terragrunt cache, this stage will be skipped and thus the hooks
      will not run. Similarly, if terragrunt detects that it does not need to run `init` in the auto init feature, the
      `init` stage is skipped along with the related hooks.
    - The working directory for hooks associated with `init-from-module` will run in the terragrunt config directory,
      while the working directory for hooks associated with `init` will be the terraform module.


Complete Example:

```hcl
terraform {
  # Pull the terraform configuration at the github repo "acme/infrastructure-modules", under the subdirectory
  # "networking/vpc", using the git tag "v0.0.1".
  source = "git::git@github.com:acme/infrastructure-modules.git//networking/vpc?ref=v0.0.1"

  # For any terraform commands that use locking, make sure to configure a lock timeout of 20 minutes.
  extra_arguments "retry_lock" {
    commands  = get_terraform_commands_that_need_locking()
    arguments = ["-lock-timeout=20m"]
  }

  # You can also specify multiple extra arguments for each use case. Here we configure terragrunt to always pass in the
  # `common.tfvars` var file located by the parent terragrunt config.
  extra_arguments "custom_vars" {
    commands = [
      "apply",
      "plan",
      "import",
      "push",
      "refresh"
    ]

    required_var_files = ["${get_parent_terragrunt_dir()}/common.tfvars"]
  }

  # The following are examples of how to specify hooks

  # Before apply or plan, run "echo Foo".
  before_hook "before_hook_1" {
    commands     = ["apply", "plan"]
    execute      = ["echo", "Foo"]
  }

  # Before apply, run "echo Bar". Note that blocks are ordered, so this hook will run after the previous hook to
  # "echo Foo". In this case, always "echo Bar" even if the previous hook failed.
  before_hook "before_hook_2" {
    commands     = ["apply"]
    execute      = ["echo", "Bar"]
    run_on_error = true
  }

  # Note that you can use interpolations in subblocks. Here, we configure it so that before apply or plan, print out the
  # environment variable "HOME".
  before_hook "interpolation_hook_1" {
    commands     = ["apply", "plan"]
    execute      = ["echo", get_env("HOME", "HelloWorld")]
    run_on_error = false
  }

  # After running apply or plan, run "echo Baz". This hook is configured so that it will always run, even if the apply
  # or plan failed.
  after_hook "after_hook_1" {
    commands     = ["apply", "plan"]
    execute      = ["echo", "Baz"]
    run_on_error = true
  }

  # After an error occurs during apply or plan, run "echo Error Hook executed". This hook is configured so that it will run
  # after any error, with the ".*" expression.
  error_hook "error_hook_1" {
    commands  = ["apply", "plan"]
    execute   = ["echo", "Error Hook executed"]
    on_errors = [
      ".*",
    ]
  }

  # A special after hook to always run after the init-from-module step of the Terragrunt pipeline. In this case, we will
  # copy the "foo.tf" file located by the parent terragrunt.hcl file to the current working directory.
  after_hook "init_from_module" {
    commands = ["init-from-module"]
    execute  = ["cp", "${get_parent_terragrunt_dir()}/foo.tf", "."]
  }

  # A special after_hook. Use this hook if you wish to run commands immediately after terragrunt finishes loading its
  # configurations. If "terragrunt-read-config" is defined as a before_hook, it will be ignored as this config would
  # not be loaded before the action is done.
  after_hook "terragrunt-read-config" {
    commands = ["terragrunt-read-config"]
    execute  = ["bash", "script/get_aws_credentials.sh"]
  }
}
```

Local File Path Example with allowed hidden files:

```hcl
terraform {
  # Pull the terraform configuration from the local file system. Terragrunt will make a copy of the source folder in the
  # Terragrunt working directory (typically `.terragrunt-cache`).
  source = "../modules/networking/vpc"

  # Always include the following file patterns in the Terragrunt copy.
  include_in_copy = [
    ".security_group_rules.json",
    "*.yaml",
  ]
}
```


#### A note about using modules from the registry

The key design of Terragrunt is to act as a preprocessor to convert **shared service modules** in the registry into a **root
module**. In Terraform, modules can be loosely categorized into two types:

* **Root Module**: A Terraform module that is designed for running `terraform init` and the other workflow commands
  (`apply`, `plan`, etc). This is the entrypoint module for deploying your infrastructure. Root modules are identified
  by the presence of key blocks that setup configuration about how Terraform behaves, like `backend` blocks (for
  configuring state) and `provider` blocks (for configuring how Terraform interacts with the cloud APIs).
* **Shared Module**: A Terraform module that is designed to be included in other Terraform modules through `module`
  blocks. These modules are missing many of the key blocks that are required for running the workflow commands of
  terraform.

Terragrunt further distinguishes shared modules between **service modules** and **modules**:

* **Shared Service Module**: A Terraform module that is designed to be standalone and applied directly. These modules
  are not root modules in that they are still missing the key blocks like `backend` and `provider`, but aside from that
  do not need any additional configuration or composition to deploy. For example, the
  [terraform-aws-modules/vpc](https://registry.terraform.io/modules/terraform-aws-modules/vpc/aws/latest) module can be
  deployed by itself without composing with other modules or resources.
* **Shared Module**: A Terraform module that is designed to be composed with other modules. That is, these modules must
  be embedded in another Terraform module and combined with other resources or modules. For example, the
  [consul-security-group-rules
  module](https://registry.terraform.io/modules/hashicorp/consul/aws/latest/submodules/consul-security-group-rules)

Terragrunt started off with features that help directly deploy **Root Modules**, but over the years have implemented
many features that allow you to turn **Shared Service Modules** into **Root Modules**  by injecting the key configuration
blocks that are necessary for Terraform modules to act as **Root Modules**.

Modules on the Terraform Registry are primarily designed to be used as **Shared Modules**. That is, you won't be able to
`git clone` the underlying repository and run `terraform init` or `apply` directly on the module without modification.
Unless otherwise specified, almost all the modules will require composition with other modules/resources to deploy.
When using modules in the registry, it helps to think about what blocks and resources are necessary to operate the
module, and translating those into Terragrunt blocks that generate them.

Note that in many cases, Terragrunt may not be able to deploy modules from the registry. While Terragrunt has features
to turn any **Shared Module** into a **Root Module**, there are two key technical limitations that prevent Terragrunt
from converting ALL shared modules:

- Every complex input must have a `type` associated with it. Otherwise, Terraform will interpret the input that
  Terragrunt passes through as `string`. This includes `list` and `map`.
- Derived sensitive outputs must be marked as `sensitive`. Refer to the [terraform tutorial on sensitive
  variables](https://learn.hashicorp.com/tutorials/terraform/sensitive-variables#reference-sensitive-variables) for more
  information on this requirement.

**If you run into issues deploying a module from the registry, chances are that module is not a Shared Service Module,
and thus not designed for use with Terragrunt. Depending on the technical limitation, Terragrunt may be able to
support the transition to root module. Please always file [an issue on the terragrunt
repository](https://github.com/gruntwork-io/terragrunt/issues) with the module + error message you are encountering,
instead of the module repository.**
