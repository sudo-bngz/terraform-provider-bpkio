package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	SlateURL     = "https://bpkiosamples.s3.eu-west-1.amazonaws.com/broadpeakio-slate.jpg"
	LiveURL      = "https://hls-radio-s3.nextradiotv.com/olyzon/delayed/master.m3u8"
	LiveURLOther = "https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8"
	AdServerURL  = "https://vast-prep.staging.olyzon.tv/sources/1042586b/serve"
)

func TestAccServiceAdInsertion_Basic(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfig(apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "id"),
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "source.id"),
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "live_ad_replacement.ad_server.id"),
					resource.TestCheckResourceAttrSet("bpkio_service_ad_insertion.test", "live_ad_replacement.gap_filler.id"),
				),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	initialName := "tf-acc-service-ad-initial"
	updatedName := "tf-acc-service-ad-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithName(apiKey, initialName),
				Check: resource.TestCheckResourceAttr(
					"bpkio_service_ad_insertion.test", "name", initialName),
			},
			{
				Config: testAccServiceAdInsertionConfigWithName(apiKey, updatedName),
				Check: resource.TestCheckResourceAttr(
					"bpkio_service_ad_insertion.test", "name", updatedName),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateSlate(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	initialName := "tf-acc-slate-initial"
	updatedName := "tf-acc-slate-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithSlateName(apiKey, initialName),
				Check:  resource.TestCheckResourceAttr("bpkio_source_slate.slate", "name", initialName),
			},
			{
				Config: testAccServiceAdInsertionConfigWithSlateName(apiKey, updatedName),
				Check:  resource.TestCheckResourceAttr("bpkio_source_slate.slate", "name", updatedName),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateLiveName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	initialName := "tf-acc-live-initial"
	updatedName := "tf-acc-live-updated"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithLiveName(apiKey, initialName),
				Check:  resource.TestCheckResourceAttr("bpkio_source_live.live", "name", initialName),
			},
			{
				Config: testAccServiceAdInsertionConfigWithLiveName(apiKey, updatedName),
				Check:  resource.TestCheckResourceAttr("bpkio_source_live.live", "name", updatedName),
			},
		},
	})
}

func TestAccServiceAdInsertion_UpdateLiveSource(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfigWithLiveSource(apiKey, LiveURL),
				Check:  resource.TestCheckResourceAttr("bpkio_source_live.live", "url", LiveURL),
			},
			{
				Config:      testAccServiceAdInsertionConfigWithLiveSource(apiKey, LiveURLOther),
				Check:       resource.TestCheckResourceAttr("bpkio_source_live.live", "url", LiveURLOther),
				ExpectError: regexp.MustCompile(`(?i)403|forbidden|not allowed`),
			},
		},
	})
}

func TestAccServiceAdInsertion_ImportStateAndDrift(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAdInsertionConfig(apiKey),
			},
			{
				ResourceName:      "bpkio_service_ad_insertion.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// // --- Error cases: try creating with invalid source/adserver id (should error out) ---

func TestAccServiceAdInsertion_InvalidSource(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	badID := 999999999
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccServiceAdInsertionConfigWithBadSource(apiKey, badID),
				ExpectError: regexp.MustCompile(`(?i)403|forbidden|not allowed`),
			},
		},
	})
}

func TestAccServiceAdInsertion_InvalidAdServer(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	badID := 999999999
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccServiceAdInsertionConfigWithBadAdServer(apiKey, badID),
				ExpectError: regexp.MustCompile(`(?i)403|forbidden|not allowed`),
			},
		},
	})
}

// --- Config helpers ---

func testAccServiceAdInsertionConfig(apiKey string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "tf-acc-slate"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "tf-acc-live"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "tf-acc-adserver"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_service_ad_insertion" "test" {
  name = "tf-acc-adinsertion"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, SlateURL, LiveURL, AdServerURL)
}

// Change live source name
func testAccServiceAdInsertionConfigWithName(apiKey, serviceName string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "tf-source-slate"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "tf-acc-live"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "tf-acc-adserver"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, SlateURL, LiveURL, AdServerURL, serviceName)
}

// Change live source name
func testAccServiceAdInsertionConfigWithLiveName(apiKey, liveName string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "tf-acc-slate"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "tf-acc-adserver"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_service_ad_insertion" "test" {
  name = "tf-acc-adinsertion"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, SlateURL, liveName, LiveURL, AdServerURL)
}

func testAccServiceAdInsertionConfigWithLiveSource(apiKey, liveURL string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "tf-acc-slate"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "tf-source-live"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "tf-acc-adserver"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_service_ad_insertion" "test" {
  name = "tf-acc-adinsertion"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, SlateURL, liveURL, AdServerURL)
}

// Change slate name
func testAccServiceAdInsertionConfigWithSlateName(apiKey, slateName string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "%s"
  url  = "%s"
}

resource "bpkio_source_live" "live" {
  name = "tf-acc-live"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "tf-acc-adserver"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_service_ad_insertion" "test" {
  name = "tf-acc-adinsertion"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, slateName, SlateURL, LiveURL, AdServerURL)
}

// Invalid source (bad id)
func testAccServiceAdInsertionConfigWithBadSource(apiKey string, badSourceID int) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_source_slate" "slate" {
  name = "tf-slate-source"
  url  = "%s"
}

resource "bpkio_source_adserver" "adserver" {
  name = "tf-acc-adserver"
  url  = "%s"
}

resource "bpkio_service_ad_insertion" "adservice" {
  name = "tf-acc-adinsertion-badsource"

  source = {
    id = %d
  }
  
  live_ad_replacement = {
    ad_server = {
      id = bpkio_source_adserver.adserver.id
    }
    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, SlateURL, AdServerURL, apiKey, badSourceID)
}

// Invalid ad server (bad id)
func testAccServiceAdInsertionConfigWithBadAdServer(apiKey string, badAdServerID int) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_live" "live" {
  name = "tf-acc-live"
  url  = "%s"
}

resource "bpkio_source_slate" "slate" {
  name = "tf-acc-live"
  url  = "%s"
}

data "bpkio_transcoding_profile" "test" {
	id = 5763
}

resource "bpkio_service_ad_insertion" "test" {
  name = "tf-acc-adinsertion-badadserver"

  source = {
    id = bpkio_source_live.live.id
  }

  live_ad_replacement = {
    ad_server = {
      id = %d
    }

    gap_filler = {
      id = bpkio_source_slate.slate.id
    }
    spot_aware = {}
  }

  transcoding_profile = {
    id = data.bpkio_transcoding_profile.test.id
  }
}
`, apiKey, LiveURL, SlateURL, badAdServerID)
}
