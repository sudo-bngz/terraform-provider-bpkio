// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sourceAdServerDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceAdServerDataSource{}
)

// sourceAdServerDataSource is the data source implementation.
type sourceAdServerDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewSourceAdServerDataSource is a helper function to simplify the provider implementation.
func NewSourceAdServerDataSource() datasource.DataSource {
	return &sourceAdServerDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *sourceAdServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *sourceAdServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_ad_server"
}

// Schema defines the schema for the data source.
func (d *sourceAdServerDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
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
			"description": schema.StringAttribute{
				Computed: true,
			},
			"queries": schema.StringAttribute{
				Computed: true,
			},
			"query_parameters": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sourceAdServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sourceAdServerDataSourceModel
	var sourceid int64

	diags := req.Config.GetAttribute(ctx, path.Root("id"), &sourceid)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Get the source from the API
	source, err := d.client.GetAdServer(uint(sourceid))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", sourceid, err.Error()),
		)
		return
	}

	sourceState := sourceAdServerDataSourceModel{
		ID:      types.Int64Value(int64(source.Id)),
		Name:    types.StringValue(source.Name),
		Type:    types.StringValue(source.Type),
		URL:     types.StringValue(source.Url),
		Queries: types.StringValue(source.Queries),
		QueryParameters: func() []queryParametersModel {
			var params []queryParametersModel
			for _, param := range source.QueryParameters {
				params = append(params, queryParametersModel{
					Type:  types.StringValue(param.Type),
					Name:  types.StringValue(param.Name),
					Value: types.StringValue(param.Value),
				})
			}
			return params
		}(),
	}

	// Set state
	diags = resp.State.Set(ctx, &sourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// sourceModel maps source schema data.
type sourceAdServerDataSourceModel struct {
	ID              types.Int64            `tfsdk:"id"`
	Name            types.String           `tfsdk:"name"`
	Description     types.String           `tfsdk:"description"`
	Type            types.String           `tfsdk:"type"`
	URL             types.String           `tfsdk:"url"`
	Queries         types.String           `tfsdk:"queries"`
	QueryParameters []queryParametersModel `tfsdk:"query_parameters"`
}
