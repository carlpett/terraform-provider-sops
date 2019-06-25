# terraform-sops

A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

## Example

Encrypt a file using Sops: `sops demo-secret.enc.json`

```json
{
  "password": "foo",
  "db": {"password": "bar"}
}
```

Usage in Terraform (0.12 and later) looks like this:

```hcl
provider "sops" {}

data "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

output "do-something" {
  value = data.sops_file.demo-secret.data["password"]
}

output "do-something2" {
  value = data.sops_file.demo-secret.data["db.password"]
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

## Install

Download the latest [release](https://github.com/carlpett/terraform-provider-sops/releases) for your environment and unpack it to the user plugin directory. The user plugins directory is in one of the following locations, depending on the host operating system:
* Windows `%APPDATA%\terraform.d\plugins`
* All other systems `~/.terraform.d/plugins`

## Development
Building and testing is most easily performed with `make build` and `make test` respectively.

The PGP key used for encrypting the test cases is found in `test/testing-key.pgp`. You can import it with `gpg --import test/testing-key.pgp`.
