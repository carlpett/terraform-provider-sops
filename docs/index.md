# sops Provider

A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

!> To prevent plaintext secrets from being written to disk, you *must* use a secure remote state backend. See the [official docs](https://www.terraform.io/docs/state/sensitive-data.html) on _Sensitive Data in State_ for more information.

## Example Usage

```hcl
provider "sops" {}

data "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

output "db-password" {
  # Access the password variable that is under db via the terraform map of data
  value = data.sops_file.demo-secret.data["db.password"]
}

data "http" "remote_sops_data" {
  url = "https://sops.example/my-data"
}

data "sops_external" "demo-secret" {
  source     = data.http.remote_sops_data.body
  input_type = "yaml"
}

output "root-value-hello" {
  value = data.sops_external.demo-secret.data.hello
}
```
