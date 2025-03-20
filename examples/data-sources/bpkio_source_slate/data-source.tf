terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_source_slate" "this" {
  id = 135320
}

output "this_source" {
  value = data.bpkio_source_slate.this
}
