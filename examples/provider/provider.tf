terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
  api_key = var.api_key
}

resource "bpkio_source_live" "this" {
  name        = "foobar-test-tf"
  description = "test"
  url         = "https://live.stream/master.m3u8"

  //TODO: Find way to handle when origin is empty
  origin = {}
}

resource "bpkio_source_adserver" "this" {
  name        = "foobar-test-tf-b"
  description = "test"
  url         = "https://ad.server/endpoint"
}
