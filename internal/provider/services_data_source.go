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
	_ datasource.DataSource              = &servicesDataSource{}
	_ datasource.DataSourceWithConfigure = &servicesDataSource{}
)

// servicesDataSource is the data source implementation.
type servicesDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewServicesDataSource is a helper function to simplify the provider implementation.
func NewServicesDataSource() datasource.DataSource {
	return &servicesDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *servicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *servicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_services"
}

// Schema defines the schema for the data source.
func (d *servicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("ad-insertion", "content-replacement", "virtual-channel"),
				},
			},
			"state": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "disabled"),
				},
			},
			"services": schema.ListNestedAttribute{
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
						"creation_date": schema.StringAttribute{
							Computed: true,
						},
						"update_date": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"tags": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *servicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state servicesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	services, err := d.client.GetAllServices(0, 2000)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read All Services",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, service := range services {
		// Convert the []string to types.List
		tagsList, diags := types.ListValueFrom(ctx, types.StringType, service.EnvironmentTags)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		serviceState := serviceDataSourceModel{
			ID:           types.Int64Value(int64(service.Id)),
			Name:         types.StringValue(service.Name),
			Type:         types.StringValue(service.Type),
			URL:          types.StringValue(service.Url),
			CreationDate: types.StringValue(service.CreationDate), // Make sure these fields exist in your API response
			UpdateDate:   types.StringValue(service.UpdateDate),
			State:        types.StringValue(service.State),
			Tags:         tagsList,
		}

		// Filter by type
		if (state.Type.IsNull() || (service.Type == state.Type.ValueString() && !state.Type.IsNull())) &&
			(state.State.IsNull() || (service.State == state.State.ValueString() && !state.State.IsNull())) {
			state.Services = append(state.Services, serviceState)
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// servicesDataSourceModel maps the data source schema data.
type servicesDataSourceModel struct {
	Type     types.String             `tfsdk:"type"`
	State    types.String             `tfsdk:"state"`
	Services []serviceDataSourceModel `tfsdk:"services"`
}

// serviceModel maps service schema data.
type serviceDataSourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	URL          types.String `tfsdk:"url"`
	CreationDate types.String `tfsdk:"creation_date"`
	UpdateDate   types.String `tfsdk:"update_date"`
	State        types.String `tfsdk:"state"`
	Tags         types.List   `tfsdk:"tags"`
}
