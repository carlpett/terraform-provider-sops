# terraform-sops

A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

**NOTE:** To prevent plaintext secrets from being written to disk, you *must* set up a secure remote state backend. See the [official docs](https://developer.hashicorp.com/terraform/language/state/sensitive-data) on _Sensitive Data in State_ for more information or use [ephemeral block](#example-using-ephemeral-block).

## Example

**NOTE:** All examples assume Terraform 0.13 or newer. For information about usage on older versions, see the [legacy usage docs](docs/legacy_usage.md).

Encrypt a file using Sops: `sops demo-secret.enc.json`

```json
{
  "password": "foo",
  "db": {"password": "bar"}
}
```
### sops_file

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}

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

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}

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

## Install

For Terraform 0.13 and later, specify the source and version in a `required_providers` block:

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}
```

## CI usage

For CI, the same variables or context that SOPS uses locally must be provided in the runtime. The provider does not manage the required values. 

## Development
Building and testing is most easily performed with `make build` and `make test` respectively.

The PGP key used for encrypting the test cases is found in `test/testing-key.pgp`. You can import it with `gpg --import test/testing-key.pgp`.

To create the Terraform-registry-documentation, simply run `make generate-documentation`

## Transitioning to Terraform 0.13 provider required blocks.

With Terraform 0.13, providers are available/downloaded via the [terraform registry](https://registry.terraform.io/providers/carlpett/sops/latest) via a required_providers block.

```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 0.5"
    }
  }
}
```

A prerequisite when converting is that you must remove the data source block from the previous SOPS provider in your `terraform.state` file. 
This can be done via:
```shell
terraform state replace-provider registry.terraform.io/-/sops registry.terraform.io/carlpett/sops
```

If not you will be greeted with: 
```shell
- Finding latest version of -/sops...

Error: Failed to query available provider packages

Could not retrieve the list of available versions for provider -/sops:
provider registry registry.terraform.io does not have a provider named
registry.terraform.io/-/sops
```

## Example using ephemeral block
With Terraform v1.11+ and the SOPS provider v1.3.0+, you can use an ephemeral resource instead of a data source.
This prevents the contents of the secret file from being saved in the Terraform state.
Ephemeral resources can be referenced in `write-only` arguments.
```hcl
terraform {
  required_providers {
    sops = {
      source = "carlpett/sops"
      version = "~> 1.3.0"
    }
  }
}

ephemeral "sops_file" "secrets" {
  source_file = "demo-secret.enc.json"
}

resource "aws_ssm_parameter" "sops_secrets" {
  name             = "my-secrets"
  type             = "SecureString"
  value_wo         = jsonencode(ephemeral.sops_file.secrets.raw)
  value_wo_version = 1
}
```
See documentation:
* [Ephemeral block](https://developer.hashicorp.com/terraform/language/block/ephemeral)
* [Write-Only arguments](https://developer.hashicorp.com/terraform/language/manage-sensitive-data/write-only)


