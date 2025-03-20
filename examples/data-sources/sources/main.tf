terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_sources" "all" {}
output "all_sources" {
  value = data.bpkio_sources.all
}

data "bpkio_sources" "slates" {
  type = "slate"
}

output "slates" {
  value = data.bpkio_sources.slates
}
