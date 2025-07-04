// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sourceLiveDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceLiveDataSource{}
)

// sourceLiveDataSource is the data source implementation.
type sourceLiveDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewSourceLiveDataSource is a helper function to simplify the provider implementation.
func NewSourceLiveDataSource() datasource.DataSource {
	return &sourceLiveDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *sourceLiveDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *sourceLiveDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_live"
}

// Schema defines the schema for the data source.
func (d *sourceLiveDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"format": schema.StringAttribute{
				Computed: true,
			},
			"multi_period": schema.BoolAttribute{
				Computed: true,
			},
			"origin": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"custom_headers": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
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
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sourceLiveDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sourceLiveDataSourceModel
	var sourceid int64

	diags := req.Config.GetAttribute(ctx, path.Root("id"), &sourceid)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Get the source from the API
	source, err := d.client.GetLive(uint(sourceid))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", sourceid, err.Error()),
		)
		return
	}

	var originAttr types.Object

	customHeaderType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	if len(source.Origin.CustomHeaders) > 0 {
		var headerValues []attr.Value

		for _, h := range source.Origin.CustomHeaders {
			headerObj, diag := types.ObjectValue(
				customHeaderType.AttrTypes,
				map[string]attr.Value{
					"name":  types.StringValue(h.Name),
					"value": types.StringValue(h.Value),
				},
			)
			if diag.HasError() {
				resp.Diagnostics.Append(diag...)
				return
			}
			headerValues = append(headerValues, headerObj)
		}

		headersList, diag := types.ListValue(customHeaderType, headerValues)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}

		originAttr, diag = types.ObjectValue(
			map[string]attr.Type{
				"custom_headers": types.ListType{ElemType: customHeaderType},
			},
			map[string]attr.Value{
				"custom_headers": headersList,
			},
		)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
	} else {
		originAttr = types.ObjectNull(map[string]attr.Type{
			"custom_headers": types.ListType{ElemType: customHeaderType},
		})
	}

	sourceState := sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		MultiPeriod: types.BoolValue(source.MultiPeriod),
		Format:      types.StringValue(source.Format),
		Origin:      originAttr,
	}

	// Set state
	diags = resp.State.Set(ctx, &sourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// sourceModel maps source schema data.
type sourceLiveDataSourceModel struct {
	ID          types.Int64           `tfsdk:"id"`
	Name        types.String          `tfsdk:"name"`
	Type        types.String          `tfsdk:"type"`
	URL         types.String          `tfsdk:"url"`
	Format      types.String          `tfsdk:"format"`
	Description types.String          `tfsdk:"description"`
	MultiPeriod types.Bool            `tfsdk:"multi_period"`
	Origin      basetypes.ObjectValue `tfsdk:"origin"`
}
