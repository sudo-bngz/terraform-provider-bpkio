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

func (d *sourceAdServerDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	//--------------------------------------------------------------------
	// 1. Parse the ID from configuration
	//--------------------------------------------------------------------
	var adServerID int64
	diags := req.Config.GetAttribute(ctx, path.Root("id"), &adServerID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//--------------------------------------------------------------------
	// 2. Call the Broadpeak API
	//--------------------------------------------------------------------
	src, err := d.client.GetAdServer(uint(adServerID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source Ad-Server",
			fmt.Sprintf("Ad-Server with ID %d not found (%s)", adServerID, err),
		)
		return
	}

	//--------------------------------------------------------------------
	// 3. Build query_parameters -> types.List
	//--------------------------------------------------------------------
	// Object schema for a single parameter
	paramObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":  types.StringType,
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	var paramValues []attr.Value
	for _, p := range src.QueryParameters {
		objVal, diag := types.ObjectValue(
			paramObjType.AttrTypes,
			map[string]attr.Value{
				"type":  types.StringValue(p.Type),
				"name":  types.StringValue(p.Name),
				"value": types.StringValue(p.Value),
			},
		)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		paramValues = append(paramValues, objVal)
	}

	var paramsList types.List
	if len(paramValues) > 0 {
		paramsList = types.ListValueMust(paramObjType, paramValues)
	} else {
		// List is absent/empty
		paramsList = types.ListNull(paramObjType)
	}

	//--------------------------------------------------------------------
	// 4. Populate Terraform state
	//--------------------------------------------------------------------
	state := sourceAdServerDataSourceModel{
		ID:              types.Int64Value(int64(src.Id)),
		Name:            types.StringValue(src.Name),
		Description:     types.StringValue(src.Description),
		Type:            types.StringValue(src.Type),
		URL:             types.StringValue(src.Url),
		Queries:         types.StringValue(src.Queries),
		QueryParameters: paramsList,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// sourceModel maps source schema data.
type sourceAdServerDataSourceModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Type            types.String `tfsdk:"type"`
	URL             types.String `tfsdk:"url"`
	Queries         types.String `tfsdk:"queries"`
	QueryParameters types.List   `tfsdk:"query_parameters"`
}
