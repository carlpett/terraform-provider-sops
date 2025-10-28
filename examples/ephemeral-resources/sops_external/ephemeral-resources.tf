data "http" "remote_sops_data" {
  url = "https://sops.example/my-data"
}

ephemeral "sops_external" "demo_secret" {
  source     = data.http.remote_sops_data.body
  input_type = "yaml"
}

output "root_value_hello" {
  value = ephemeral.sops_external.demo_secret.data.hello
}