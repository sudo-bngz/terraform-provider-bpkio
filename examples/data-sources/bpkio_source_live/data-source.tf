terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_source_live" "this" {
  id = 123082
}

output "this_source" {
  value = data.bpkio_source_live.this
}
