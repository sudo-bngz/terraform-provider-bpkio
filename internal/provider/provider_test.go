// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bpkio-terraform-provider/internal/testacc"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"bpkio": func() (tfprotov6.ProviderServer, error) {
			p := New("dev")()
			factory := providerserver.NewProtocol6WithError(p)
			return factory()
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("BPKIO_ENDPOINT"); v == "" {
		t.Fatal("BPKIO_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("BPKIO_API_KEY"); v == "" {
		t.Fatal("BPKIO_API_KEY must be set for acceptance tests")
	}
}

/* ------------------------------------------------------------------------- */
/* Acceptance test                                                            */
/* ------------------------------------------------------------------------- */
func TestAccProvider_basic(t *testing.T) {

	mock := testacc.NewMock()
	defer mock.Close()

	// --- make env-vars visible to the provider ----
	os.Setenv("BPKIO_ENDPOINT", mock.URL())
	os.Setenv("BPKIO_API_KEY", mock.Token())
	defer func() {
		os.Unsetenv("BPKIO_ENDPOINT")
		os.Unsetenv("BPKIO_API_KEY")
	}()

	/* 3) Run the Terraform test case */
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProviderFactories(),
		PreCheck:                 func() { testAccPreCheck(t) },

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "bpkio" {
  api_key  = "%s"
  endpoint = "%s"
}
`, "mock.Token()", mock.URL()),
			},
		},
	})
}
