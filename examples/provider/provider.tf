terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
  api_key = "changeme"
}
