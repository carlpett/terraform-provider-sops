# Read only the sops metadata - the file is never decrypted, so no data key
# access is required and no secret values are written to state.
data "sops_file_metadata" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

# Read the secret itself via an ephemeral resource so it stays out of state.
ephemeral "sops_file" "demo-secret" {
  source_file = "demo-secret.enc.json"
}

# The secret is write-only (never stored in state); the persisted
# last_modified_unix drives the version, so Terraform re-pushes the value
# whenever the encrypted file is re-encrypted.
resource "aws_ssm_parameter" "demo-secret" {
  name             = "demo-secret"
  type             = "SecureString"
  value_wo         = ephemeral.sops_file.demo-secret.data["password"]
  value_wo_version = data.sops_file_metadata.demo-secret.last_modified_unix
}
