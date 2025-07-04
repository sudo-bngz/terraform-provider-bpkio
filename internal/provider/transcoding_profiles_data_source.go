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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// --------------------------------------------------------------------
// Interface assertions
// --------------------------------------------------------------------
var (
	_ datasource.DataSource              = &transcodingProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &transcodingProfilesDataSource{}
)

// --------------------------------------------------------------------
// Data-source definition
// --------------------------------------------------------------------
type transcodingProfilesDataSource struct {
	client *broadpeakio.BroadpeakClient
}

func NewTranscodingProfilesDataSource() datasource.DataSource {
	return &transcodingProfilesDataSource{}
}

// --------------------------------------------------------------------
// Configure
// --------------------------------------------------------------------
func (d *transcodingProfilesDataSource) Configure(
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
			"Unexpected Provider Data",
			fmt.Sprintf("Expected *broadpeakio.BroadpeakClient, got %T", req.ProviderData),
		)
		return
	}
	d.client = c
}

// --------------------------------------------------------------------
// Metadata
// --------------------------------------------------------------------
func (d *transcodingProfilesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_profiles"
}

// --------------------------------------------------------------------
// Schema
// --------------------------------------------------------------------
func (d *transcodingProfilesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"profiles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":          schema.Int64Attribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"content":     schema.StringAttribute{Computed: true},
						"internal_id": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

// --------------------------------------------------------------------
// Read
// --------------------------------------------------------------------
func (d *transcodingProfilesDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// 1. Call the Broadpeak API
	list, err := d.client.GetAllTranscodingProfiles(0, 2000)
	if err != nil {
		resp.Diagnostics.AddError("Unable to List Transcoding Profiles", err.Error())
		return
	}

	// 2. Build Terraform-typed list
	profileObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.Int64Type,
			"name":        types.StringType,
			"content":     types.StringType,
			"internal_id": types.StringType,
		},
	}

	var objs []attr.Value
	for _, p := range list {
		objVal, diag := types.ObjectValue(
			profileObjType.AttrTypes,
			map[string]attr.Value{
				"id":          types.Int64Value(int64(p.Id)),
				"name":        types.StringValue(p.Name),
				"content":     types.StringValue(string(p.Content)), // Raw JSON → string
				"internal_id": types.StringValue(p.InternalId),
			},
		)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		objs = append(objs, objVal)
	}

	var profilesList types.List
	if len(objs) > 0 {
		profilesList = types.ListValueMust(profileObjType, objs)
	} else {
		profilesList = types.ListNull(profileObjType)
	}

	// 3. Set state
	state := transcodingProfilesDataSourceModel{
		Profiles: profilesList,
	}
	diag := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diag...)
}

// --------------------------------------------------------------------
// State model
// --------------------------------------------------------------------
type transcodingProfilesDataSourceModel struct {
	Profiles types.List `tfsdk:"profiles"` // List<Object>
}
