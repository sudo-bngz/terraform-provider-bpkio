terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_transcoding_profiles" "all" {}

output "all_transcoding_profiles" {
  value = data.bpkio_transcoding_profiles.all
}
