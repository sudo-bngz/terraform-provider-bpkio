package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Valid Ad server
var adServerURL = "https://vast-prep.staging.olyzon.tv/sources/1042586b/serve"

func TestAccSourceAdServer_Basic(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	if apiKey == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
	resourceName := "bpkio_source_adserver.test"

	// This URL MUST point to a real, working adserver in your environment!
	adServerURL := "https://vast-prep.staging.olyzon.tv/sources/1042586b/serve"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceAdServerConfig(apiKey, adServerURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-test-adserver"),
					resource.TestCheckResourceAttr(resourceName, "url", adServerURL),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
			{
				// Test update: Change description and check result
				Config: testAccSourceAdServerConfigUpdate(apiKey, adServerURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated by acceptance test"),
				),
			},
		},
	})
}

func testAccSourceAdServerConfig(apiKey, url string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  name = "tf-acc-test-adserver"
  url  = "%s"
}
`, apiKey, url)
}

func testAccSourceAdServerConfigUpdate(apiKey, url string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  name        = "tf-acc-test-adserver"
  url         = "%s"
  description = "updated by acceptance test"
}
`, apiKey, url)
}

func TestAccSourceAdServer_ComputedFields(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	resourceName := "bpkio_source_adserver.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceAdServerConfig(apiKey, adServerURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				),
			},
		},
	})
}

func TestAccSourceAdServer_MinimalConfig(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	resourceName := "bpkio_source_adserver.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceAdServerConfig(apiKey, adServerURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-test-adserver"),
					resource.TestCheckResourceAttr(resourceName, "url", adServerURL),
				),
			},
		},
	})
}

func TestAccSourceAdServer_MissingName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  url = "%s"
}
`, apiKey, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)The argument "name" is required, but no definition was found.`),
			},
		},
	})
}

func TestAccSourceAdServer_LongSpecialName(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	longName := strings.Repeat("x", 101) // >100 chars to trigger validation error
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}

resource "bpkio_source_adserver" "test" {
  name = "%s"
  url  = "%s"
}
`, apiKey, longName, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`(?s)Bad Request`),
			},
		},
	})
}

func TestAccSourceAdServer_DuplicateNameURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	unique := "tf-acc-dupe"
	dupeConfig := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}
resource "bpkio_source_adserver" "first" {
  name = "%s"
  url  = "%s"
}
resource "bpkio_source_adserver" "second" {
  name = "%s"
  url  = "%s"
}
`, apiKey, unique, adServerURL, unique, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      dupeConfig,
				ExpectError: regexp.MustCompile(`(?s)403|500`),
			},
		},
	})
}

func TestAccSourceAdServer_QueryParameters(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	resourceName := "bpkio_source_adserver.test"
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}
resource "bpkio_source_adserver" "test" {
  name = "tf-acc-test-adserver"
  url  = "%s"
  query_parameters = [
    {
      type  = "from-header"
      name  = "X-Test"
      value = "value1"
    },
    {
      type  = "custom"
      name  = "X-Custom"
      value = "value2"
    }
  ]
}
`, apiKey, adServerURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "query_parameters.#", "2"),
				),
			},
		},
	})
}

func TestAccSourceAdServer_InvalidURL(t *testing.T) {
	apiKey := os.Getenv("BPKIO_API_KEY")
	badURL := "https://this-url-will-not-exist.example.com"
	config := fmt.Sprintf(`
provider "bpkio" {
  api_key = "%s"
}
resource "bpkio_source_adserver" "test" {
  name = "tf-acc-test-adserver"
  url  = "%s"
}
`, apiKey, badURL)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
		},
	})
}
