terraform {
  required_providers {
    bpkio = {
      source = "bashou/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_transcoding_profiles" "all" {}

output "all_transcoding_profiles" {
  value = data.bpkio_transcoding_profiles.all
}
