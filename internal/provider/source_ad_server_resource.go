// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strconv"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

func (r *sourceAdServerResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	//--------------------------------------------------------------------
	// 1. Decode the plan into a strongly-typed model
	//--------------------------------------------------------------------
	var plan sourceAdServerDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//--------------------------------------------------------------------
	// 2. Build the Broadpeak API input
	//--------------------------------------------------------------------
	adInput := broadpeakio.AdServerInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Url:         plan.URL.ValueString(),
		Queries:     plan.Queries.ValueString(),
	}

	// Decode query_parameters list → slice for API
	if !plan.QueryParameters.IsNull() && !plan.QueryParameters.IsUnknown() {
		var paramSlice []queryParametersModel
		diags := plan.QueryParameters.ElementsAs(ctx, &paramSlice, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, p := range paramSlice {
			adInput.QueryParameters = append(adInput.QueryParameters, broadpeakio.QueryParam{
				Type:  p.Type.ValueString(),
				Name:  p.Name.ValueString(),
				Value: p.Value.ValueString(),
			})
		}
	}

	//--------------------------------------------------------------------
	// 3. Call Broadpeak to create the Ad-Server
	//--------------------------------------------------------------------
	created, err := r.client.CreateAdServer(adInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Ad-Server",
			fmt.Sprintf("Could not create Ad-Server: %s", err),
		)
		return
	}

	//--------------------------------------------------------------------
	// 4. Convert API → Terraform state
	//--------------------------------------------------------------------
	// Object type for a single parameter
	paramObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":  types.StringType,
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	var paramValues []attr.Value
	for _, p := range created.QueryParameters {
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
		paramsList = types.ListNull(paramObjType)
	}

	newState := sourceAdServerDataSourceModel{
		ID:              types.Int64Value(int64(created.Id)),
		Name:            types.StringValue(created.Name),
		Description:     types.StringValue(created.Description),
		Type:            types.StringValue(created.Type),
		URL:             types.StringValue(created.Url),
		Queries:         types.StringValue(created.Queries),
		QueryParameters: paramsList,
	}

	//--------------------------------------------------------------------
	// 5. Save state
	//--------------------------------------------------------------------
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sourceAdServerResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	//--------------------------------------------------------------------
	// 1. Load the prior state (contains the ID)
	//--------------------------------------------------------------------
	var state sourceAdServerDataSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//--------------------------------------------------------------------
	// 2. Query Broadpeak for the latest object
	//--------------------------------------------------------------------
	src, err := r.client.GetAdServer(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Source Ad-Server",
			fmt.Sprintf("Ad-Server with ID %d not found (%s)", state.ID.ValueInt64(), err),
		)
		return
	}

	//--------------------------------------------------------------------
	// 3. Convert QueryParameters -> types.List
	//--------------------------------------------------------------------
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
		paramsList = types.ListNull(paramObjType)
	}

	//--------------------------------------------------------------------
	// 4. Build the new state object
	//--------------------------------------------------------------------
	newState := sourceAdServerDataSourceModel{
		ID:              types.Int64Value(int64(src.Id)),
		Name:            types.StringValue(src.Name),
		Description:     types.StringValue(src.Description),
		Type:            types.StringValue(src.Type),
		URL:             types.StringValue(src.Url),
		Queries:         types.StringValue(src.Queries),
		QueryParameters: paramsList,
	}

	//--------------------------------------------------------------------
	// 5. Save state
	//--------------------------------------------------------------------
	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (r *sourceAdServerResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	//--------------------------------------------------------------------
	// 1. Decode the planned values
	//--------------------------------------------------------------------
	var plan sourceAdServerDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//--------------------------------------------------------------------
	// 2. Build the Broadpeak input
	//--------------------------------------------------------------------
	updInput := broadpeakio.AdServerInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Url:         plan.URL.ValueString(),
		Queries:     plan.Queries.ValueString(),
		Template:    "custom",
	}

	if !plan.QueryParameters.IsNull() && !plan.QueryParameters.IsUnknown() {
		var paramSlice []queryParametersModel
		diags := plan.QueryParameters.ElementsAs(ctx, &paramSlice, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for _, p := range paramSlice {
			updInput.QueryParameters = append(updInput.QueryParameters, broadpeakio.QueryParam{
				Type:  p.Type.ValueString(),
				Name:  p.Name.ValueString(),
				Value: p.Value.ValueString(),
			})
		}
	}

	//--------------------------------------------------------------------
	// 3. Call the Broadpeak API
	//--------------------------------------------------------------------
	adID := uint(plan.ID.ValueInt64())
	if _, err := r.client.UpdateAdServer(adID, updInput); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Ad-Server",
			fmt.Sprintf("Could not update ad-server ID %d: %s", adID, err),
		)
		return
	}

	//--------------------------------------------------------------------
	// 4. Re-query to obtain the authoritative object
	//--------------------------------------------------------------------
	src, err := r.client.GetAdServer(adID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Ad-Server",
			fmt.Sprintf("Could not fetch ad-server ID %d after update: %s", adID, err),
		)
		return
	}

	//--------------------------------------------------------------------
	// 5. Build query_parameters -> types.List
	//--------------------------------------------------------------------
	paramObjType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":  types.StringType,
			"name":  types.StringType,
			"value": types.StringType,
		},
	}

	var paramVals []attr.Value
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
		paramVals = append(paramVals, objVal)
	}

	var paramsList types.List
	if len(paramVals) > 0 {
		paramsList = types.ListValueMust(paramObjType, paramVals)
	} else {
		paramsList = types.ListNull(paramObjType)
	}

	//--------------------------------------------------------------------
	// 6. Write the new state
	//--------------------------------------------------------------------
	newState := sourceAdServerDataSourceModel{
		ID:              types.Int64Value(int64(src.Id)),
		Name:            types.StringValue(src.Name),
		Description:     types.StringValue(src.Description),
		Type:            types.StringValue(src.Type),
		URL:             types.StringValue(src.Url),
		Queries:         types.StringValue(src.Queries),
		QueryParameters: paramsList,
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
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
