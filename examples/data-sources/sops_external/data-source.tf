data "http" "remote_sops_data" {
  url = "https://sops.example/my-data"
}

data "sops_external" "demo_secret" {
  source     = data.http.remote_sops_data.body
  input_type = "yaml"
}

output "root-value-hello" {
  value = data.sops_external.demo_secret.data.hello
}