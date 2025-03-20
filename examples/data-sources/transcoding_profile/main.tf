terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_transcoding_profile" "this" {
  id = 4694
}

output "this_transcoding_profile" {
  value = data.bpkio_transcoding_profile.this
}
