terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_source_ad_server" "this" {
  id = 132658
}

output "this_source" {
  value = data.bpkio_source_ad_server.this
}
