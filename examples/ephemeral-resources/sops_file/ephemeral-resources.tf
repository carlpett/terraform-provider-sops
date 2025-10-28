provider "sops" {}

ephemeral "sops_file" "demo_secret" {
  source_file = "demo_secret.enc.json"
}

output "root_value_password" {
  # Access the password variable from the map
  value = data.sops_file.demo_secret.data["password"]
}

output "mapped_nested_value" {
  # Access the password variable that is under db via the terraform map of data
  value = data.sops_file.demo_secret.data["db.password"]
}

output "nested_json_value" {
  # Access the password variable that is under db via the terraform object
  value = jsondecode(data.sops_file.demo_secret.raw).db.password
}