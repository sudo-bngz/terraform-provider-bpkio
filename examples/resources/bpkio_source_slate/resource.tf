terraform {
  required_providers {
    bpkio = {
      source = "bashou/bpkio"
    }
  }
}

provider "bpkio" {
}

resource "bpkio_source_slate" "this" {
  name        = "foobar-test-tf"
  description = "test slate"
  url         = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4"
}
