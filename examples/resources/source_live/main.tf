terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}

resource "bpkio_source_live" "this" {
  name        = "foobar-test-tf"
  description = "test"
  url         = "https://hls-radio-s3.nextradiotv.com/lg/bfmtv/master.m3u8"

  //TODO: Find way to handle when origin is empty
  origin = {}
}
