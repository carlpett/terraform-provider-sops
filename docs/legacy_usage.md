# terraform-sops on older Terraform versions
## Migrating existing states
To migrate a state from Terraform 0.12 or older, there is a need to change how the provider is referenced. Terraform provides a command to do this migration:

```shell
terraform state replace-provider registry.terraform.io/-/sops registry.terraform.io/carlpett/sops
```

## Installation

Download the latest [release](https://github.com/carlpett/terraform-provider-sops/releases) for your environment and unpack it to the user plugin directory. The user plugins directory is in one of the following locations, depending on the host operating system:
* Windows `%APPDATA%\terraform.d\plugins`
* All other systems `~/.terraform.d/plugins`

### Allowing code to run on macOS

Apple macOS Catalina (10.15.0) and later prevents unsigned code from running. When you first run `terraform plan` it will pop up a message saying
> **“terraform-provider-sops_v0.5.0” cannot be opened because the developer cannot be verified.**
> macOS cannot verify that this app is free from malware.

To allow the plugin to run, go to the **Security & Privacy** tab of System Preferences and you should see a message saying
> “terraform-provider-sops_v0.5.0” was blocked from use because it is not from an identified developer.

Click the `Allow Anyway` button.

## Usage
Usage is mostly identical across versions, but there are some differences in how to reference nested fields.

### Terraform 0.12

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

### Terraform 0.11 and older
```hcl
provider "sops" {}

data "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

output "do-something" {
  value = "${data.sops_file.demo-secret.data.password}"
}

output "do-something2" {
  value = "${data.sops_file.demo-secret.data.db.password}"
}
```
