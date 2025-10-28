provider "sops" {}

data "sops_file" "demo_secret" {
  source_file = "demo-secret.enc.json"
}

output "db_password" {
  # Access the password variable that is under db via the terraform map of data
  value = data.sops_file.demo_secret.data["db.password"]
}
