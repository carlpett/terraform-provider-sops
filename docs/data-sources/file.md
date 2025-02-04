# sops_file Data Source

Read data from a sops-encrypted file on disk.

## Example Usage

```hcl
provider "sops" {}

data "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

output "root-value-password" {
  # Access the password variable from the map
  value = data.sops_file.demo-secret.data["password"]
}

output "mapped-nested-value" {
  # Access the password variable that is under db via the terraform map of data
  value = data.sops_file.demo-secret.data["db.password"]
}

output "nested-json-value" {
  # Access the password variable that is under db via the terraform object
  value = jsondecode(data.sops_file.demo-secret.raw).db.password
}
```

## Argument Reference

* `source_file` - (Required) Path to the encrypted file
* `input_type` - (Optional) The provider will use the file extension to determine how to unmarshal the data. If your file does not have the usual extension, set this argument to `yaml`, `json`, `dotenv` (`.env`), `ini` accordingly, or `raw` if the encrypted data is encoded differently.

## Attribute Reference

* `data` - The unmarshalled data as a dictionary. Use dot-separated keys to access nested data.
* `raw` - The entire unencrypted file as a string.
