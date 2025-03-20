terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

#TODO: Create resource directly and remove hardcoded ID
data "bpkio_service_ad_insertion" "this" {
  id = 54235
}

output "service_output" {
  value = data.bpkio_service_ad_insertion.this
}
