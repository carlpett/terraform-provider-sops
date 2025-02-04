# sops_external Data Source

Read data from a sops-encrypted string. Useful if the data does not reside on disk locally (otherwise use `sops_file`).

## Example Usage

```hcl
provider "sops" {}

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

## Argument Reference

* `source` - (Required) A string with sops-encrypted data
* `input_type` - (Required) `yaml`, `json` `dotenv` (`.env`), `ini` or `raw`, depending on the structure of the un-encrypted data.

## Attribute Reference

* `data` - The unmarshalled data as a dictionary. Use dot-separated keys to access nested data.
* `raw` - The entire unencrypted file as a string.
