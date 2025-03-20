// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &bpkioProvider{}
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &bpkioProvider{
			version: version,
		}
	}
}

// bpkioProvider is the provider implementation.
type bpkioProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *bpkioProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "bpkio"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *bpkioProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Required:    true,
				Description: "API key for Broadpeak",
				Sensitive:   true,
			},
		},
	}
}

// bpkioProviderModel maps provider schema data to a Go type.
type bpkioProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

// Configure prepares a bpkio API client for data sources and resources.
func (p *bpkioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config bpkioProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown bpkio API Endpoint",
			"The provider cannot create the bpkio API client as there is an unknown configuration value for the bpkio API Endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BPKIO_ENDPOINT environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown bpkio API Key",
			"The provider cannot create the bpkio API client as there is an unknown configuration value for the bpkio API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BPKIO_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	endpoint := getenv("BPKIO_ENDPOINT", "https://api.broadpeak.io")
	api_key := getenv("BPKIO_API_KEY", "")

	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}

	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	tflog.Debug(ctx, endpoint)
	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing bpkio API Key",
			"The provider cannot create the bpkio API client as there is a missing or empty value for the bpkio API key. "+
				"Set the username value in the configuration or use the BPKIO_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new bpkio client using the configuration values
	//TODO: Find a way to test key
	client := broadpeakio.MakeClient(api_key)
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Unable to Create bpkio API Client",
	//		"An unexpected error occurred when creating the bpkio API client. "+
	//			"If the error is not clear, please contact the provider developers.\n\n"+
	//			"bpkio Client Error: "+err.Error(),
	//	)
	//	return
	//}

	// Make the bpkio client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &client
	resp.ResourceData = &client
}

// DataSources defines the data sources implemented in the provider.
func (p *bpkioProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSourcesDataSource,
		NewSourceAdServerDataSource,
		NewSourceSlateDataSource,
		NewSourceLiveDataSource,
		NewServiceAdInsertionDataSource,
		NewServicesDataSource,
		NewTranscodingProfileDataSource,
		NewTranscodingProfilesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *bpkioProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewServiceAdInsertionResource,
		NewSourceSlateResource,
		NewSourceLiveResource,
		NewSourceAdServerResource,
	}
}
