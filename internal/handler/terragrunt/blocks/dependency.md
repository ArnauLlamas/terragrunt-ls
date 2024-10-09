dependency
---
The `dependency` block is used to configure module dependencies. Each dependency block exports the outputs of the target
module as block attributes you can reference throughout the configuration. You can learn more about `dependency` blocks
in the [Dependencies between modules
section](/docs/features/execute-terraform-commands-on-multiple-modules-at-once/#dependencies-between-modules) of the
"Execute Terraform commands on multiple modules at once" use case overview.

You can define more than one `dependency` block. Each label you provide to the block identifies another `dependency`
that you can reference in your config.

The `dependency` block supports the following arguments:

- `name` (label): You can define multiple `dependency` blocks in a single terragrunt config. As such, each block needs a
  name to differentiate between the other blocks, which is what the first label of the block is used for. You can
  reference the specific dependency output by the name. E.g if you had a block `dependency "vpc"`, you can reference the
  outputs of this dependency with the expression `dependency.vpc.outputs`.
- `config_path` (attribute): Path to a Terragrunt module (folder with a `terragrunt.hcl` file) that should be included
  as a dependency in this configuration.
- `enabled` (attribute): When `false`, excludes the dependency from execution. Defaults to `true`.
- `skip_outputs` (attribute): When `true`, skip calling `terragrunt output` when processing this dependency. If
  `mock_outputs` is configured, set `outputs` to the value of `mock_outputs`. Otherwise, `outputs` will be set to an
  empty map. Put another way, setting `skip_outputs` means "use mocks all the time if `mock_outputs` are set."
- `mock_outputs` (attribute): A map of arbitrary key value pairs to use as the `outputs` attribute when no outputs are
  available from the target module, or if `skip_outputs` is `true`. However, it's generally recommended not to set
  `skip_outputs` if using `mock_outputs`, because `skip_outputs` means "use mocks all the time if they are set" whereas
  `mock_outputs` means "use mocks only if real outputs are not available." Use `locals` instead when `skip_outputs = true`.
- `mock_outputs_allowed_terraform_commands` (attribute): A list of Terraform commands for which `mock_outputs` are
  allowed. If a command is used where `mock_outputs` is not allowed, and no outputs are available in the target module,
  Terragrunt will throw an error when processing this dependency.
- `mock_outputs_merge_with_state` (attribute): DEPRECATED. Use `mock_outputs_merge_strategy_with_state`. When `true`,
  `mock_outputs` and the state outputs will be merged. That is, the `mock_outputs` will be treated as defaults and the
  real state outputs will overwrite them if the keys clash.
- `mock_outputs_merge_strategy_with_state` (attribute): Specifies how any existing state should be merged into the
  mocks. Valid values are
  - `no_merge` (default) - any existing state will be used as is. If the dependency does not have an existing state (it
    hasn't been applied yet), then the mocks will be used
  - `shallow` - the existing state will be shallow merged into the mocks. Mocks will only be used where the output does
    not already exist in the dependency's state
  - `deep_map_only` - the existing state will be deeply merged into the mocks. If an output is a map, the mock key
    will be used where that key does not exist in the state. Lists will not be merged

Example:

```hcl
# Run `terragrunt output` on the module at the relative path `../vpc` and expose them under the attribute
# `dependency.vpc.outputs`
dependency "vpc" {
  config_path = "../vpc"

  # Configure mock outputs for the `validate` command that are returned when there are no outputs available (e.g the
  # module hasn't been applied yet.
  mock_outputs_allowed_terraform_commands = ["validate"]
  mock_outputs = {
    vpc_id = "fake-vpc-id"
  }
}

# Another dependency, available under the attribute `dependency.rds.outputs`
dependency "rds" {
  config_path = "../rds"
}

inputs = {
  vpc_id = dependency.vpc.outputs.vpc_id
  db_url = dependency.rds.outputs.db_url
}
```

**Can I speed up dependency fetching?**

`dependency` blocks are fetched in parallel at each source level, but will serially parse each recursive dependency. For
example, consider the following chain of dependencies:

```
account --> vpc --> securitygroup --> ecs
                                      ^
                                     /
                              ecr --
```

In this chain, the `ecr` and `securitygroup` module outputs will be fetched concurrently when applying the `ecs` module,
but the outputs for `account` and `vpc` will be fetched serially as terragrunt needs to recursively walk through the
tree to retrieve the outputs at each level.

This recursive parsing happens due to the necessity to parse the entire `terragrunt.hcl` configuration (including
`dependency` blocks) in full before being able to call `terraform output`.

However, terragrunt includes an optimization to only fetch the lowest level outputs (`securitygroup` and `ecr` in this
example) provided that the following conditions are met in the immediate dependencies:

- The remote state is managed using `remote_state` blocks.
- The dependency optimization feature flag is enabled (`disable_dependency_optimization = false`, which is the default).
- The `remote_state` block itself does not depend on any `dependency` outputs (`locals` and `include` are ok).
- You are not relying on `before_hook`, `after_hook`, or `extra_arguments` to the `terraform init` call. NOTE:
  terragrunt will not automatically detect this and you will need to explicitly opt out of the dependency optimization
  flag.

If these conditions are met, terragrunt will only parse out the `remote_state` blocks and use that to pull down the
state for the target module without parsing the `dependency` blocks, avoiding the recursive dependency retrieval.
