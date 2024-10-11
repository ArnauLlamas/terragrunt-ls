locals {
  a  = "1"
  aa = read_terragrunt_config(find_in_parent_folders("root.hcl"))
}

inputs = {}

dependency "pepe" {}
