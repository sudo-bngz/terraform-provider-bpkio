terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
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
