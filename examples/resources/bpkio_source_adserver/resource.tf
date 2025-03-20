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
  url         = "https://hls-radio-s3.nextradiotv.com/lg/bfmtv/master.m3u8"

  //TODO: Handle case when query_parameters is empty
  query_parameters = []
}
