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
  url         = "https://live.stream/master.m3u8"

  //TODO: Find way to handle when origin is empty
  origin = {}
}

resource "bpkio_source_slate" "this" {
  name        = "foobar-test-tf"
  description = "test slate"
  url         = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4"
}
