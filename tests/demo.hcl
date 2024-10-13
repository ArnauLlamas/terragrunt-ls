locals {
  one  = "1"
  file = read_terragrunt_config(find_in_parent_folders("root.hcl"))
}

include "root" {
  path = "root.hcl"
}

include "patch" {
  path = "patch.hcl"
}

dependency "name" {}

inputs = {

}
