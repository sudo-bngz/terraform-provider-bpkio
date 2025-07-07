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

// --------------------------------------------------------------------
// Type assertions
// --------------------------------------------------------------------
var (
	_ datasource.DataSource              = &transcodingProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &transcodingProfileDataSource{}
)

// --------------------------------------------------------------------
// Data-source definition
// --------------------------------------------------------------------
type transcodingProfileDataSource struct {
	client *broadpeakio.BroadpeakClient
}

func NewTranscodingProfileDataSource() datasource.DataSource {
	return &transcodingProfileDataSource{}
}

// --------------------------------------------------------------------
// Configure
// --------------------------------------------------------------------
func (d *transcodingProfileDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*broadpeakio.BroadpeakClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *broadpeakio.BroadpeakClient, got %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

// --------------------------------------------------------------------
// Metadata
// --------------------------------------------------------------------
func (d *transcodingProfileDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_profile"
}

// --------------------------------------------------------------------
// Schema
// --------------------------------------------------------------------
func (d *transcodingProfileDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			// We keep raw JSON as a string for simplicity
			"content": schema.StringAttribute{
				Computed: true,
			},
			"internal_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// --------------------------------------------------------------------
// Read
// --------------------------------------------------------------------
func (d *transcodingProfileDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Parse ID from config
	var id int64
	diags := req.Config.GetAttribute(ctx, path.Root("id"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch profile from API
	p, err := d.client.GetTranscodingProfile(uint(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Transcoding Profile",
			fmt.Sprintf("Profile ID %d not found: %s", id, err),
		)
		return
	}

	// Build state
	state := transcodingProfileDataSourceModel{
		ID:         types.Int64Value(int64(p.Id)),
		Name:       types.StringValue(p.Name),
		Content:    types.StringValue(string(p.Content)),
		InternalId: types.StringValue(p.InternalId),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// --------------------------------------------------------------------
// State model
// --------------------------------------------------------------------
type transcodingProfileDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Content    types.String `tfsdk:"content"`
	InternalId types.String `tfsdk:"internal_id"`
}
