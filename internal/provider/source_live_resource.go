// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sourceLiveResource{}
	_ resource.ResourceWithConfigure   = &sourceLiveResource{}
	_ resource.ResourceWithImportState = &sourceLiveResource{}
)

// NewSourceLiveResource is a helper function to simplify the provider implementation.
func NewSourceLiveResource() resource.Resource {
	return &sourceLiveResource{}
}

// sourceLiveResource is the resource implementation.
type sourceLiveResource struct {
	client *broadpeakio.BroadpeakClient
}

// Configure adds the provider configured client to the resource.
func (r *sourceLiveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Metadata returns the resource type name.
func (r *sourceLiveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_live"
}

// Schema defines the schema for the resource.
func (r *sourceLiveResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
			},
			"format": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"multi_period": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"origin": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"custom_headers": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required: true,
								},
								"value": schema.StringAttribute{
									Required: true,
								},
							},
						},
					},
				},
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *sourceLiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sourceLiveDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform model to API model
	var sourceData = broadpeakio.LiveInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		MultiPeriod: plan.MultiPeriod.ValueBool(),
		Url:         plan.URL.ValueString(),
	}

	// Handle Origin data if present
	if plan.Origin != nil {
		var headers []broadpeakio.CustomHeader
		for _, header := range plan.Origin.CustomHeaders {
			headers = append(headers, broadpeakio.CustomHeader{
				Name:  header.Name.ValueString(),
				Value: header.Value.ValueString(),
			})
		}
		sourceData.Origin = broadpeakio.Origin{
			CustomHeaders: headers,
		}
	}

	// Create new live
	source, err := r.client.CreateLive(sourceData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating live",
			"Could not create live, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	result := sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Format:      types.StringValue(source.Format),
		MultiPeriod: types.BoolValue(source.MultiPeriod),

		Origin: &originModel{
			CustomHeaders: func() []customHeadersModel {
				var headers []customHeadersModel
				for _, header := range source.Origin.CustomHeaders {
					headers = append(headers, customHeadersModel{
						Name:  types.StringValue(header.Name),
						Value: types.StringValue(header.Value),
					})
				}
				return headers
			}(),
		},
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *sourceLiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state sourceLiveDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed live value from HashiCups
	source, err := r.client.GetLive(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", state.ID, err.Error()),
		)
		return
	}

	state = sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Format:      types.StringValue(source.Format),
		MultiPeriod: types.BoolValue(source.MultiPeriod),

		Origin: &originModel{
			CustomHeaders: func() []customHeadersModel {
				var headers []customHeadersModel
				for _, header := range source.Origin.CustomHeaders {
					headers = append(headers, customHeadersModel{
						Name:  types.StringValue(header.Name),
						Value: types.StringValue(header.Value),
					})
				}
				return headers
			}(),
		},
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sourceLiveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan and current state
	var plan sourceLiveDataSourceModel

	// Get planned changes
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the update data
	var sourceData = broadpeakio.LiveInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		MultiPeriod: plan.MultiPeriod.ValueBool(),
		Url:         plan.URL.ValueString(),

		Origin: broadpeakio.Origin{
			CustomHeaders: func() []broadpeakio.CustomHeader {
				var headers []broadpeakio.CustomHeader
				for _, header := range plan.Origin.CustomHeaders {
					headers = append(headers, broadpeakio.CustomHeader{
						Name:  header.Name.ValueString(),
						Value: header.Value.ValueString(),
					})
				}
				return headers
			}(),
		},
	}

	// Retrieve ID from plan/state
	liveID := uint(plan.ID.ValueInt64())
	// Update existing live
	_, err := r.client.UpdateLive(liveID, sourceData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating live",
			"Could not update live, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetLive
	source, err := r.client.GetLive(liveID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Live",
			fmt.Sprintf("Could not fetch live ID %d: %s", liveID, err.Error()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	result := sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Format:      types.StringValue(source.Format),
		MultiPeriod: types.BoolValue(source.MultiPeriod),

		Origin: &originModel{
			CustomHeaders: func() []customHeadersModel {
				var headers []customHeadersModel
				for _, header := range source.Origin.CustomHeaders {
					headers = append(headers, customHeadersModel{
						Name:  types.StringValue(header.Name),
						Value: types.StringValue(header.Value),
					})
				}
				return headers
			}(),
		},
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sourceLiveResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sourceLiveDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing live
	_, err := r.client.DeleteLive(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source Live",
			"Could not delete live, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state from the ID.
func (r *sourceLiveResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Convert the ID from string to int64
	idStr := req.ID

	// Parse the ID string into an int
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing source live",
			fmt.Sprintf("Invalid ID format: %s. Expected a numeric ID. Error: %s", idStr, err),
		)
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)

	// After importing the ID, the Read method will be called automatically to refresh the state
}
