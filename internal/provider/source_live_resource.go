// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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
				Computed:    true,
				Description: "The ID of the source live.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the source live.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the source live.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL of the source live.",
			},
			"format": schema.StringAttribute{
				Computed:    true,
				Description: "The format of the source live.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the source live.",
				Default:     stringdefault.StaticString(""),
			},
			"multi_period": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the source live supports multiple periods.(Default: `false`)",
				Default:     booldefault.StaticBool(false),
			},
			"origin": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"custom_headers": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Required:    true,
									Description: "The name of the custom header.",
								},
								"value": schema.StringAttribute{
									Required:    true,
									Description: "The value of the custom header.",
								},
							},
						},
					},
				},
				Optional:    true,
				Computed:    true,
				Description: "The origin configuration for the source live.",
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *sourceLiveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve the plan into a strongly typed model
	var plan sourceLiveDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the API input from the Terraform plan
	sourceData := broadpeakio.LiveInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		MultiPeriod: plan.MultiPeriod.ValueBool(),
		Url:         plan.URL.ValueString(),
	}

	// Handle optional origin block
	if !plan.Origin.IsNull() && !plan.Origin.IsUnknown() {
		var origin originModel
		diags := plan.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, h := range origin.CustomHeaders {
			sourceData.Origin.CustomHeaders = append(sourceData.Origin.CustomHeaders, broadpeakio.CustomHeader{
				Name:  h.Name.ValueString(),
				Value: h.Value.ValueString(),
			})
		}
	}

	// Call the Broadpeak API to create the resource
	source, err := r.client.CreateLive(sourceData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating source live",
			fmt.Sprintf("Could not create source live: %s", err),
		)
		return
	}

	// Build origin attribute for Terraform state
	customHeaderType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	originAttrType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"custom_headers": types.ListType{ElemType: customHeaderType},
		},
	}

	var originAttr types.Object
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

		originAttr, diag = types.ObjectValue(originAttrType.AttrTypes, map[string]attr.Value{
			"custom_headers": headersList,
		})
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
	} else {
		originAttr = types.ObjectNull(originAttrType.AttrTypes)
	}

	// Build the final Terraform state model
	result := sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Format:      types.StringValue(source.Format),
		MultiPeriod: types.BoolValue(source.MultiPeriod),
		Origin:      originAttr,
	}

	// Save the state
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sourceLiveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceLiveDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	source, err := r.client.GetLive(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source Live",
			fmt.Sprintf("Service with ID %d not found (%s)", state.ID.ValueInt64(), err.Error()),
		)
		return
	}

	// Build the origin object attribute
	customHeaderType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	var headerValues []attr.Value
	for _, header := range source.Origin.CustomHeaders {
		headerVal, diag := types.ObjectValue(
			customHeaderType.AttrTypes,
			map[string]attr.Value{
				"name":  types.StringValue(header.Name),
				"value": types.StringValue(header.Value),
			},
		)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		headerValues = append(headerValues, headerVal)
	}

	originAttrType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"custom_headers": types.ListType{ElemType: customHeaderType},
		},
	}

	var originAttr types.Object
	if len(headerValues) > 0 {
		headersList := types.ListValueMust(types.ObjectType{AttrTypes: customHeaderType.AttrTypes}, headerValues)
		originVal, diag := types.ObjectValue(
			originAttrType.AttrTypes,
			map[string]attr.Value{
				"custom_headers": headersList,
			},
		)
		if diag.HasError() {
			resp.Diagnostics.Append(diag...)
			return
		}
		originAttr = originVal
	} else {
		originAttr = types.ObjectNull(originAttrType.AttrTypes)
	}

	// Set state
	state = sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Format:      types.StringValue(source.Format),
		MultiPeriod: types.BoolValue(source.MultiPeriod),
		Origin:      originAttr,
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sourceLiveResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// ---------------------------------------------------------------------
	// 1. Load the planned state
	// ---------------------------------------------------------------------
	var plan sourceLiveDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// ---------------------------------------------------------------------
	// 2. Build LiveInput for the Broadpeak API
	// ---------------------------------------------------------------------
	updateInput := broadpeakio.LiveInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		MultiPeriod: plan.MultiPeriod.ValueBool(),
		Url:         plan.URL.ValueString(),
	}

	// Decode origin from plan if present
	if !plan.Origin.IsNull() && !plan.Origin.IsUnknown() {
		var origin originModel
		diags := plan.Origin.As(ctx, &origin, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only include headers if at least one exists
		for _, h := range origin.CustomHeaders {
			updateInput.Origin.CustomHeaders = append(
				updateInput.Origin.CustomHeaders,
				broadpeakio.CustomHeader{
					Name:  h.Name.ValueString(),
					Value: h.Value.ValueString(),
				},
			)
		}
	}

	// ---------------------------------------------------------------------
	// 3. Call the API to update
	// ---------------------------------------------------------------------
	liveID := uint(plan.ID.ValueInt64())
	if _, err := r.client.UpdateLive(liveID, updateInput); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Source Live",
			fmt.Sprintf("Could not update source live ID %d: %s", liveID, err),
		)
		return
	}

	// ---------------------------------------------------------------------
	// 4. Re-query the updated object so the state is authoritative
	// ---------------------------------------------------------------------
	source, err := r.client.GetLive(liveID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Source Live",
			fmt.Sprintf("Could not retrieve source live ID %d after update: %s", liveID, err),
		)
		return
	}

	// ---------------------------------------------------------------------
	// 5. Convert origin from API -> types.Object for Terraform
	// ---------------------------------------------------------------------
	customHeaderType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		},
	}
	originType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"custom_headers": types.ListType{ElemType: customHeaderType},
		},
	}

	var originAttr types.Object
	if len(source.Origin.CustomHeaders) > 0 {
		var headerVals []attr.Value
		for _, h := range source.Origin.CustomHeaders {
			hv, diag := types.ObjectValue(
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
			headerVals = append(headerVals, hv)
		}

		listVal := types.ListValueMust(customHeaderType, headerVals)

		var diag diag.Diagnostics
		originAttr, diag = types.ObjectValue(originType.AttrTypes, map[string]attr.Value{
			"custom_headers": listVal,
		})
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		originAttr = types.ObjectNull(originType.AttrTypes)
	}

	// ---------------------------------------------------------------------
	// 6. Write final state
	// ---------------------------------------------------------------------
	newState := sourceLiveDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Format:      types.StringValue(source.Format),
		MultiPeriod: types.BoolValue(source.MultiPeriod),
		Origin:      originAttr,
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
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
