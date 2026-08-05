provider "sops" {}

# data_wo accepts any Terraform value directly, preserving real types instead of jsonencode/yamlencode
resource "sops_file" "demo-secret" {
  destination_file = "demo-secret.enc.yaml"
  data_wo = {
    password = "correct-horse-battery-staple"
    db = {
      password = "hunter2"
    }
    port    = 5432
    enabled = true
  }
  data_wo_version = "1"
}
