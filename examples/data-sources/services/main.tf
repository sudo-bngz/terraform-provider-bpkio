terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

data "bpkio_services" "all" {}

output "all_services" {
  value = data.bpkio_services.all
}

data "bpkio_services" "ad-insertion" {
  type = "ad-insertion"
}

output "ad-insertion" {
  value = data.bpkio_services.ad-insertion
}

data "bpkio_services" "content-replacement" {
  type = "content-replacement"
}

output "content-replacement" {
  value = data.bpkio_services.content-replacement
}

data "bpkio_services" "disabled" {
  state = "disabled"
}

output "disabled" {
  value = data.bpkio_services.disabled
}
