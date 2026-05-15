data "http" "remote_sops_data" {
  url = "https://sops.example/my-data"
}

ephemeral "sops_external" "demo_secret" {
  source     = data.http.remote_sops_data.body
  input_type = "yaml"
}

resource "sops_hash" "demo_secret" {
  input_wo = ephemeral.sops_external.demo_secret.data
}

resource "time_static" "demo_secret_update" {
  triggers = {
    value = sops_hash.demo_secret.hash
  }
}

resource "aws_secretsmanager_secret_version" "demo_secret" {
  secret_id                = "demo_secret"
  secret_string_wo         = ephemeral.sops_external.demo_secret.data
  secret_string_wo_version = time_static.demo_secret_update.unix
}