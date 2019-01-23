# terraform-sops

A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

## Example

Encrypt a file using Sops: `sops demo-secret.enc.json`

``` json
{
  "password": "foo",
  "db": {"password": "bar"}
}
```

Usage in Terraform looks like this:

``` hcl
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

## Install

``` shell
go get github.com/carlpett/terraform-sops
mkdir -p ~/.terraform.d/plugins
ln -s $GOPATH/bin/terraform-sops $HOME/.terraform.d/plugins/terraform-provider-sops
```

## Development
Building and testing is most easily performed with `make build` and `make test` respectively.

The PGP key used for encrypting the test cases is found in `test/testing-key.pgp`. You can import it with `gpg --import test/testing-key.pgp`.
