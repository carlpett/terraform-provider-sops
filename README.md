# terraform-sops

A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

**NOTE:** To prevent plaintext secrets from being written to disk, you *must* set up a secure remote state backend. See the [official docs](https://www.terraform.io/docs/state/sensitive-data.html) on _Sensitive Data in State_ for more information.

## Example

Encrypt a file using Sops: `sops demo-secret.enc.json`

```json
{
  "password": "foo",
  "db": {"password": "bar"}
}
```
### sops_file
Usage in Terraform (0.12 and later) looks like this:

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

<details><summary><i>Expand for older, Terraform 0.11 and earlier, syntax</i></summary>

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
</details>

Sops also supports encrypting the entire file when in other formats. Such files can also be used by specifying `input_type = "raw"`:

```hcl
data "sops_file" "some-file" {
  source_file = "secret-data.txt"
  input_type = "raw"
}

output "do-something" {
  value = data.sops_file.some-file.raw
}
```

### sops_external
For use with reading files that might not be local. 

> `input_type` is required with this data source.

Terraform 0.12
```hcl
provider "sops" {}

# using sops/test-fixtures/basic.yaml as an example
data "local_file" "yaml" {
  filename = "basic.yaml"
}

data "sops_external" "demo-secret" {
  source     = data.local_file.yaml.content
  input_type = "yaml"
}

output "root-value-hello" {
  value = data.sops_external.demo-secret.data.hello
}

output "nested-yaml-value" {
  # Access the password variable that is under db via the terraform object
  value = yamldecode(data.sops_file.demo-secret.raw).db.password
}
```

<details><summary><i>Expand for older, Terraform 0.11 and earlier, syntax</i></summary>

> `input_type` is required with this data source.

```hcl
provider "sops" {}

# using sops/test-fixtures/basic.yaml as an example
data "local_file" "yaml" {
  filename = "basic.yaml"
}

data "sops_external" "demo-secret" {
  source     = "${data.local_file.yaml.content}"
  input_type = "yaml"
}

output "do-something" {
  value = "${data.sops_external.demo-secret.data.hello}"
}
```
</details>



## Install

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

## Development
Building and testing is most easily performed with `make build` and `make test` respectively.

The PGP key used for encrypting the test cases is found in `test/testing-key.pgp`. You can import it with `gpg --import test/testing-key.pgp`.
