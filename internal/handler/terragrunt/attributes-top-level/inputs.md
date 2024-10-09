inputs
---
The `inputs` attribute is a map that is used to specify the input variables and their values to pass in to Terraform.
Each entry of the map will be passed to Terraform using [the environment variable
mechanism](https://www.terraform.io/docs/configuration/variables.html#environment-variables). This means that each input
will be set using the form `TF_VAR_variablename`, with the value in `json` encoded format.

Note that because the values are being passed in with environment variables and `json`, the type information is lost
when crossing the boundary between Terragrunt and Terraform. You must specify the proper [type
constraint](https://www.terraform.io/docs/configuration/variables.html#type-constraints) on the variable in Terraform in
order for Terraform to process the inputs to the right type.

Example:

```hcl
inputs = {
  string      = "string"
  number      = 42
  bool        = true
  list_string = ["a", "b", "c"]
  list_number = [1, 2, 3]
  list_bool   = [true, false]

  map_string = {
    foo = "bar"
  }

  map_number = {
    foo = 42
    bar = 12345
  }

  map_bool = {
    foo = true
    bar = false
    baz = true
  }

  object = {
    str  = "string"
    num  = 42
    list = [1, 2, 3]

    map = {
      foo = "bar"
    }
  }

  from_env = get_env("FROM_ENV", "default")
}
```
