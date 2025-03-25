// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sourceAdServerResource{}
	_ resource.ResourceWithConfigure   = &sourceAdServerResource{}
	_ resource.ResourceWithImportState = &sourceAdServerResource{}
)

// NewSourceAdServerResource is a helper function to simplify the provider implementation.
func NewSourceAdServerResource() resource.Resource {
	return &sourceAdServerResource{}
}

// sourceAdServerResource is the resource implementation.
type sourceAdServerResource struct {
	client *broadpeakio.BroadpeakClient
}

// Configure adds the provider configured client to the resource.
func (r *sourceAdServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *sourceAdServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_adserver"
}

// Schema defines the schema for the resource.
func (r *sourceAdServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the adserver. This is a unique identifier for the adserver resource.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the adserver.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the adserver. This is a read-only field.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL of the adserver.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the adserver. This field is optional and can be used to provide additional information about the adserver.",
				Default:     stringdefault.StaticString(""),
			},
			"queries": schema.StringAttribute{
				Computed:           true,
				Optional:           true,
				DeprecationMessage: "This field is deprecated and will be removed in future versions. Use 'query_parameters' instead.",
				Description:        "The queries associated with the adserver. This field is optional and can be used to specify additional query parameters for the adserver.",
				Default:            stringdefault.StaticString(""),
			},
			"query_parameters": schema.ListNestedAttribute{
				Computed:    true,
				Optional:    true,
				Description: "A list of query parameters for the adserver. Each parameter has a type, name, and value.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of the query parameter. This field is required and must be one of the following values: 'from-query-parameter', 'from-variable', 'from-header', 'forward', or 'custom'.",
							Validators: []validator.String{
								stringvalidator.OneOf("from-query-parameter", "from-variable", "from-header", "forward", "custom"),
							},
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of the query parameter. This field is required and must be a valid string.",
						},
						"value": schema.StringAttribute{
							Required:    true,
							Description: "The value of the query parameter. This field is required and must be a valid string.",
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *sourceAdServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sourceAdServerDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform model to API model
	var sourceData = broadpeakio.AdServerInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Url:         plan.URL.ValueString(),
		Queries:     plan.Queries.ValueString(),
	}

	// Handle QueryParameters data if present
	if len(plan.QueryParameters) > 0 {
		// Process parameters
		var parameters []broadpeakio.QueryParam
		for _, param := range plan.QueryParameters {
			parameters = append(parameters, broadpeakio.QueryParam{
				Type:  param.Type.ValueString(),
				Name:  param.Name.ValueString(),
				Value: param.Value.ValueString(),
			})
		}
		sourceData.QueryParameters = parameters
	} else {
		// Ensure we use empty slice, not null
		sourceData.QueryParameters = []broadpeakio.QueryParam{}
	}

	// Create new adserver
	source, err := r.client.CreateAdServer(sourceData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating adserver",
			"Could not create adserver, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	result := sourceAdServerDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Queries:     types.StringValue(source.Queries),
	}

	// Handle QueryParameters data if present
	if len(source.QueryParameters) > 0 {
		// Process parameters
		var params []queryParametersModel
		for _, param := range source.QueryParameters {
			params = append(params, queryParametersModel{
				Type:  types.StringValue(param.Type),
				Name:  types.StringValue(param.Name),
				Value: types.StringValue(param.Value),
			})
		}

		result.QueryParameters = params
	} else {
		// Ensure we use empty slice, not null
		result.QueryParameters = []queryParametersModel{}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *sourceAdServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state sourceAdServerDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed adserver value from HashiCups
	source, err := r.client.GetAdServer(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", state.ID, err.Error()),
		)
		return
	}

	state = sourceAdServerDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Queries:     types.StringValue(source.Queries),
	}

	// Handle QueryParameters data if present
	if len(source.QueryParameters) > 0 {
		// Process parameters
		var params []queryParametersModel
		for _, param := range source.QueryParameters {
			params = append(params, queryParametersModel{
				Type:  types.StringValue(param.Type),
				Name:  types.StringValue(param.Name),
				Value: types.StringValue(param.Value),
			})
		}

		state.QueryParameters = params
	} else {
		// Ensure we use empty slice, not null
		state.QueryParameters = []queryParametersModel{}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sourceAdServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan and current state
	var plan sourceAdServerDataSourceModel
	// Get planned changes
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the update data
	var sourceData = broadpeakio.AdServerInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Url:         plan.URL.ValueString(),
		Queries:     plan.Queries.ValueString(),
	}

	// Handle QueryParameters data if present
	if len(plan.QueryParameters) > 0 {
		// Process parameters
		var parameters []broadpeakio.QueryParam
		for _, param := range plan.QueryParameters {
			parameters = append(parameters, broadpeakio.QueryParam{
				Type:  param.Type.ValueString(),
				Name:  param.Name.ValueString(),
				Value: param.Value.ValueString(),
			})
		}
		sourceData.QueryParameters = parameters
	} else {
		// Ensure we use empty slice, not null
		sourceData.QueryParameters = []broadpeakio.QueryParam{}
	}

	// Retrieve ID from plan/state
	adserverID := uint(plan.ID.ValueInt64())
	// Update existing adserver
	_, err := r.client.UpdateAdServer(adserverID, sourceData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating adserver",
			"Could not update adserver, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetAdServer
	source, err := r.client.GetAdServer(adserverID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AdServer",
			fmt.Sprintf("Could not fetch adserver ID %d: %s", adserverID, err.Error()),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	result := sourceAdServerDataSourceModel{
		ID:          types.Int64Value(int64(source.Id)),
		Name:        types.StringValue(source.Name),
		Type:        types.StringValue(source.Type),
		URL:         types.StringValue(source.Url),
		Description: types.StringValue(source.Description),
		Queries:     types.StringValue(source.Queries),
	}

	// Handle QueryParameters data if present
	if len(source.QueryParameters) > 0 {
		// Process parameters
		var params []queryParametersModel
		for _, param := range source.QueryParameters {
			params = append(params, queryParametersModel{
				Type:  types.StringValue(param.Type),
				Name:  types.StringValue(param.Name),
				Value: types.StringValue(param.Value),
			})
		}

		result.QueryParameters = params
	} else {
		// Ensure we use empty slice, not null
		result.QueryParameters = []queryParametersModel{}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sourceAdServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sourceAdServerDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing adserver
	_, err := r.client.DeleteAdServer(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source AdServer",
			"Could not delete adserver, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state from the ID.
func (r *sourceAdServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Convert the ID from string to int64
	idStr := req.ID

	// Parse the ID string into an int
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing source adserver",
			fmt.Sprintf("Invalid ID format: %s. Expected a numeric ID. Error: %s", idStr, err),
		)
		return
	}

	// Set the ID in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)

	// After importing the ID, the Read method will be called automatically to refresh the state
}
