// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProvider_basic(t *testing.T) {
	//   TF_ACC=1 go test ./... -v
	if testing.Short() {
		t.Skip("skipping acceptance test in short mode")
	}

	resource.Test(t, resource.TestCase{
		// Framework → protocole v6
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"bpkio": providerserver.NewProtocol6WithError(
				New("dev")(),
			),
		},

		Steps: []resource.TestStep{
			{
				Config: `
					terraform {
					  required_providers {
					    bpkio = {
					      source  = "registry.terraform.io/bashou/bpkio"
					      version = "0.0.0"
					    }
					  }
					}

					provider "bpkio" {
					  api_key = "dummy-token"
					  // Si votre provider lit des variables d’environnement
					  // (ex : BPKIO_ENDPOINT), ajoutez-les via t.Setenv()
					}
				`,
			},
		},
	})
}
