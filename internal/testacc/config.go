// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package testacc

import (
	"fmt"
	"os"
)

func AdInsertionConfig(name string) string {
	return fmt.Sprintf(`
provider "bpkio" {
  api_key = "dummy-token"
  endpoint = "%s"
}

resource "bpkio_service_ad_insertion" "test" {
  name = "%s"
}
`, mockEndpointEnv(), name)
}

/*
mockEndpointEnv() just fetches the env var that the provider’s Configure()
uses to override the base URL; implement as you like, e.g.:
*/
func mockEndpointEnv() string {
	return os.Getenv("BPKIO_ENDPOINT")
}
