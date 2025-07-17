// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sourcesDataSource{}
	_ datasource.DataSourceWithConfigure = &sourcesDataSource{}
)

// sourcesDataSource is the data source implementation.
type sourcesDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewSourcesDataSource is a helper function to simplify the provider implementation.
func NewSourcesDataSource() datasource.DataSource {
	return &sourcesDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *sourcesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*broadpeakio.BroadpeakClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *broadpeakio.BroadpeakClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *sourcesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sources"
}

// Schema defines the schema for the data source.
func (d *sourcesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("live", "asset", "asset-catalog", "slate", "ad-server"),
				},
			},
			"sources": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"url": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sourcesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sourcesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	sources, err := d.client.GetAllSources(0, 2000)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Sources",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, source := range sources {
		sourceState := sourcesModel{
			ID:   types.Int64Value(int64(source.Id)),
			Name: types.StringValue(source.Name),
			Type: types.StringValue(source.Type),
			URL:  types.StringValue(source.Url),
		}

		if state.Type.IsNull() || (source.Type == state.Type.ValueString() && !state.Type.IsNull()) {
			state.Sources = append(state.Sources, sourceState)
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func flattenSources(sdkSources []broadpeakio.Source, filterType *string) []sourcesModel {
	var result []sourcesModel

	for _, s := range sdkSources {
		if filterType != nil && *filterType != "" && s.Type != *filterType {
			continue
		}
		result = append(result, sourcesModel{
			ID:   types.Int64Value(int64(s.Id)),
			Name: types.StringValue(s.Name),
			Type: types.StringValue(s.Type),
			URL:  types.StringValue(s.Url),
		})
	}

	if result == nil {
		return []sourcesModel{} // 👈 force non-nil return
	}

	return result
}

// sourcesDataSourceModel maps the data source schema data.
type sourcesDataSourceModel struct {
	Type    types.String   `tfsdk:"type"`
	Sources []sourcesModel `tfsdk:"sources"`
}

// sourcesModel maps sources schema data.
type sourcesModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	URL  types.String `tfsdk:"url"`
}
