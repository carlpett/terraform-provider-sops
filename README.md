# terraform-sops
A Terraform plugin for using files encrypted with [Mozilla sops](https://github.com/mozilla/sops).

# Example
Encrypt a file using Sops: `sops demo-secret.enc.json`

```json
{
	"password": "foo"
}
```

Usage in Terraform looks like this:
```hcl
provider "sops" {}

data "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

output "do-something" {
  value = "${data.sops_file.demo-secret.data.password}"
}
```
