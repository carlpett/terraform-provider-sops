# Read only the sops metadata from an encrypted string, without decrypting it.
data "sops_external_metadata" "demo-secret" {
  source     = file("demo-secret.enc.yaml")
  input_type = "yaml"
}

output "demo-secret-last-modified" {
  value = data.sops_external_metadata.demo-secret.last_modified
}
