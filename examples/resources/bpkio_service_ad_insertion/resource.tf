terraform {
  required_providers {
    bpkio = {
      source = "rmcbfm.io/terraform/bpkio"
    }
  }
}

provider "bpkio" {
}


resource "bpkio_source_adserver" "this" {
  name        = "foobar-test-tf-adserver"
  description = "test"
  url         = "https://hls-radio-s3.nextradiotv.com/lg/bfmtv/master.m3u8"

  //TODO: Handle case when query_parameters is empty
  query_parameters = []
}

resource "bpkio_source_live" "this" {
  name        = "foobar-test-tf-liveorigin"
  description = "test"
  url         = "https://hls-radio-s3.nextradiotv.com/lg/bfmtv/master.m3u8"

  //TODO: Find way to handle when origin is empty
  origin = {}
}

resource "bpkio_source_slate" "this" {
  name        = "foobar-test-tf-slate"
  description = "test slate"
  url         = "http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4"
}

data "bpkio_transcoding_profile" "this" {
  id = 4694
}

resource "bpkio_service_ad_insertion" "this" {
  name = "foobar-test-tf"

  live_ad_preroll = {
    ad_server = {
      id = bpkio_source_adserver.this.id
    }
    max_duration = 10
  }
  server_side_ad_tracking = {}

  source = {
    id = bpkio_source_live.this.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.this.id
    }

    gap_filler = {
      id = bpkio_source_slate.this.id
    }

    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.this.id
  }
}


resource "bpkio_service_ad_insertion" "this_no_preroll" {
  name = "foobar-test-tf-no-preroll-bis"

  source = {
    id = bpkio_source_live.this.id
  }

  server_side_ad_tracking = {}

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.this.id
    }

    gap_filler = {
      id = bpkio_source_slate.this.id
    }

    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.this.id
  }
}

# ID 53826
resource "bpkio_service_ad_insertion" "test_free" {
  name = "[TEST TERRAFORM] BFM2 - Production - FAI FREE - FR - back(eu-west-3)"

  advanced_options = {
    authorization_header = {
      name  = "X-BPKIO-TOKEN"
      value = "3826aad4e408cb7c3918ce63620d3b56"
    }
  }

  source = {
    id = 120450
  }

  server_side_ad_tracking = {}

  live_ad_replacement = {
    ad_server = {
      id = 125334
    }

    gap_filler = {
      id = 120491
    }

    spot_aware = {
      mode = "disabled"
    }
  }

  transcoding_profile = {
    id = 5007
  }

  tags = [
    "production",
    "bfm2",
    "fai",
    "free",
    "video",
  ]
}


data "bpkio_service_ad_insertion" "this" {
  id = bpkio_service_ad_insertion.this.id
}
output "service_ad_insertion_id" {
  value = bpkio_service_ad_insertion.this.id
}


data "bpkio_service_ad_insertion" "this_no_preroll" {
  id = bpkio_service_ad_insertion.this_no_preroll.id
}
output "service_adinsertion_no_preroll_id" {
  value = bpkio_service_ad_insertion.this_no_preroll.id
}


data "bpkio_service_ad_insertion" "test_free" {
  id = bpkio_service_ad_insertion.test_free.id
}
output "service_adinsertion_test_free_source" {
  value = bpkio_service_ad_insertion.test_free.source
}
