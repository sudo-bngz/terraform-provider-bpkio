terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
}

resource "bpkio_source_adserver" "this" {
  name        = "foobar-test-tf-b"
  description = "test"
  url         = "https://ad.server/endpoint"

  //TODO: Handle case when query_parameters is empty
  query_parameters = []
}
