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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serviceAdInsertionResource{}
	_ resource.ResourceWithConfigure   = &serviceAdInsertionResource{}
	_ resource.ResourceWithImportState = &serviceAdInsertionResource{}
)

// NewServiceAdInsertionResource is a helper function to simplify the provider implementation.
func NewServiceAdInsertionResource() resource.Resource {
	return &serviceAdInsertionResource{}
}

// serviceAdInsertionResource is the resource implementation.
type serviceAdInsertionResource struct {
	client *broadpeakio.BroadpeakClient
}

// Configure adds the provider configured client to the resource.
func (r *serviceAdInsertionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *serviceAdInsertionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_ad_insertion"
}

// Schema defines the schema for the resource.
func (r *serviceAdInsertionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Ad Insertion service creation (see https://developers.broadpeak.io/reference/adinsertioncontroller_create_v1).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "ID of the ad insertion service. This is a unique identifier for the service.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the ad insertion service. This is a human-readable name for the service.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of the ad insertion service. This indicates the type of service being created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Computed:    true,
				Description: "URL of the ad insertion service. This is the endpoint where the service can be accessed.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date of the ad insertion service. This indicates when the service was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"update_date": schema.StringAttribute{
				WriteOnly:   true,
				Optional:    true,
				Description: "Update date of the ad insertion service. This indicates when the service was last updated.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "State of the ad insertion service. This indicates the current state of the service. Possible values are 'enabled', 'paused', or 'bypassed'.",
				Default:     stringdefault.StaticString("enabled"),
				Validators: []validator.String{
					stringvalidator.OneOf("enabled", "paused", "bypassed"),
				},
			},
			"tags": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Tags for the ad insertion service. This is a list of tags associated with the service.",
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"live_ad_replacement": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"ad_server": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Required:    true,
								Description: "ID of the ad server. This is a unique identifier for the ad server.",
							},
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "Name of the ad server. This is a human-readable name for the ad server.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"type": schema.StringAttribute{
								Computed:    true,
								Description: "Type of the ad server. This indicates the type of ad server being used.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"url": schema.StringAttribute{
								Computed:    true,
								Description: "URL of the ad server. This is the endpoint where the ad server can be accessed.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional:    true,
						Description: "Ad server configuration. This is the ad server used for ad replacement.",
					},
					"gap_filler": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Required:    true,
								Description: "ID of the gap filler. This is a unique identifier for the gap filler.",
							},
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "Name of the gap filler. This is a human-readable name for the gap filler.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"type": schema.StringAttribute{
								Computed:    true,
								Description: "Type of the gap filler. This indicates the type of gap filler being used.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"url": schema.StringAttribute{
								Computed:    true,
								Description: "URL of the gap filler. This is the endpoint where the gap filler can be accessed.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional:    true,
						Description: "Gap filler configuration. This is the gap filler / slate used for ad replacement.",
					},
					"spot_aware": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Default: stringdefault.StaticString("disabled"),
								Validators: []validator.String{
									stringvalidator.OneOf("french_addressable_tv", "spot_to_live", "disabled"),
								},
								Computed:    true,
								Optional:    true,
								Description: "Mode of the spot aware. This indicates the mode of the spot aware feature. (valid values are 'french_addressable_tv', 'spot_to_live', or 'disabled'. Default: `disabled`).",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional: true,
					},
				},
				Optional:    true,
				Description: "Live ad replacement configuration. This is the configuration for live ad replacement.",
			},
			"live_ad_preroll": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"ad_server": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"id": schema.Int64Attribute{
								Required:    true,
								Description: "ID of the ad server. This is a unique identifier for the ad server.",
							},
							"name": schema.StringAttribute{
								Computed:    true,
								Description: "Name of the ad server. This is a human-readable name for the ad server.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"type": schema.StringAttribute{
								Computed:    true,
								Description: "Type of the ad server. This indicates the type of ad server being used.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"url": schema.StringAttribute{
								Computed:    true,
								Description: "URL of the ad server. This is the endpoint where the ad server can be accessed.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Optional:    true,
						Description: "Ad server configuration. This is the ad server used for ad pre-roll.",
					},
					"max_duration": schema.Int64Attribute{
						Optional:    true,
						Description: "Pre-roll maximum duration (in seconds)",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"offset": schema.Int64Attribute{
						Computed:    true,
						Optional:    true,
						Description: "Pre-roll relative start time (in seconds)",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
				Optional: true,
			},
			"advanced_options": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"authorization_header": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"value": schema.StringAttribute{
								Optional: true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
						Sensitive: true,
						Optional:  true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_ad_transcoding": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Default: booldefault.StaticBool(true),
			},
			"server_side_ad_tracking": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable": schema.BoolAttribute{
						Computed: true,
						Optional: true,
						Default:  booldefault.StaticBool(true),
					},
					"check_ad_media_segment_availability": schema.BoolAttribute{
						Computed: true,
						Optional: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
						Default: booldefault.StaticBool(false),
					},
				},
				Optional: true,
			},
			"source": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Required: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"type": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"url": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"format": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"description": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"multi_period": schema.BoolAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Optional: true,
			},
			"transcoding_profile": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Required: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"internal_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"content": schema.StringAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Optional: true,
			},
		},
	}

}

// Create creates the resource and sets the initial Terraform state.
func (r *serviceAdInsertionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan serviceAdInsertionResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform model to API model
	var tags []string
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var serviceData = broadpeakio.CreateAdInsertionInput{
		Name: plan.Name.ValueString(),
		Tags: tags,
	}

	// Add TranscodingProfile if provided
	if plan.TranscodingProfile != nil {
		serviceData.TranscodingProfile = &broadpeakio.Identifiable{
			Id: uint(plan.TranscodingProfile.ID.ValueInt64()),
		}
	}

	// Add Source if provided
	if plan.Source != nil {
		serviceData.Source = &broadpeakio.Identifiable{
			Id: uint(plan.Source.ID.ValueInt64()),
		}
	}

	// Add EnableAdTranscoding if provided
	if !plan.EnableAdTranscoding.IsNull() {
		serviceData.EnableAdTranscoding = plan.EnableAdTranscoding.ValueBool()
	}

	// Add LiveAdPreRoll if provided
	if plan.LiveAdPreRoll != nil {
		serviceData.LiveAdPreRoll = &broadpeakio.LiveAdPreRoll{
			MaxDuration: uint(plan.LiveAdPreRoll.MaxDuration.ValueInt64()),
			Offset:      uint(plan.LiveAdPreRoll.Offset.ValueInt64()),
		}

		// Add AdServer if provided within LiveAdPreRoll
		if !plan.LiveAdPreRoll.AdServer.ID.IsNull() {
			serviceData.LiveAdPreRoll.AdServer = &broadpeakio.Identifiable{
				Id: uint(plan.LiveAdPreRoll.AdServer.ID.ValueInt64()),
			}
		}
	}

	// Add LiveAdReplacement if provided
	if plan.LiveAdReplacement != nil {
		serviceData.LiveAdReplacement = &broadpeakio.LiveAdReplacement{}

		// Add AdServer if provided within LiveAdReplacement
		if !plan.LiveAdReplacement.AdServer.ID.IsNull() {
			serviceData.LiveAdReplacement.AdServer = &broadpeakio.Identifiable{
				Id: uint(plan.LiveAdReplacement.AdServer.ID.ValueInt64()),
			}
		}

		// Add GapFiller if provided
		if !plan.LiveAdReplacement.GapFiller.ID.IsNull() {
			serviceData.LiveAdReplacement.GapFiller = &broadpeakio.Identifiable{
				Id: uint(plan.LiveAdReplacement.GapFiller.ID.ValueInt64()),
			}
		}

		// Add SpotAware if provided
		if !plan.LiveAdReplacement.SpotAware.Mode.IsNull() {
			serviceData.LiveAdReplacement.SpotAware = broadpeakio.SpotAware{
				Mode: plan.LiveAdReplacement.SpotAware.Mode.ValueString(),
			}
		}
	}

	// Add ServerSideAdTracking if provided
	if plan.ServerSideAdTracking != nil {
		serviceData.ServerSideAdTracking = &broadpeakio.ServerSideAdTracking{
			Enable:                          plan.ServerSideAdTracking.Enable.ValueBool(),
			CheckAdMediaSegmentAvailability: plan.ServerSideAdTracking.CheckAdMediaSegmentAvailability.ValueBool(),
		}
	}

	if plan.AdvancedOptions != nil && plan.AdvancedOptions.AuthorizationHeader != nil {
		serviceData.AdvancedOptions = &broadpeakio.AdvancedOptions{
			AuthorizationHeader: broadpeakio.AuthorizationHeader{
				Name:  plan.AdvancedOptions.AuthorizationHeader.Name.ValueString(),
				Value: plan.AdvancedOptions.AuthorizationHeader.Value.ValueString(),
			},
		}
	}

	// Create new adserver
	service, err := r.client.CreateAdInsertion(serviceData)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating adserver",
			"Could not create adserver, unexpected error: "+err.Error(),
		)
		return
	}

	// Convert the []string to types.List
	tagsList, diags := types.ListValueFrom(ctx, types.StringType, service.Tags)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	// Map response body to schema and populate Computed attribute values
	result := serviceAdInsertionResourceModel{
		ID:                  types.Int64Value(int64(service.Id)),
		Name:                types.StringValue(service.Name),
		Type:                types.StringValue(service.Type),
		URL:                 types.StringValue(service.Url),
		CreationDate:        types.StringValue(service.CreationDate),
		UpdateDate:          types.StringValue(service.UpdateDate),
		State:               types.StringValue(service.State),
		Tags:                tagsList,
		EnableAdTranscoding: types.BoolValue(service.EnableAdTranscoding),
	}

	// Add ServerSideAdTracking if provided
	if service.ServerSideAdTracking.Enable {
		result.ServerSideAdTracking = &serverSideAdTrackingModel{
			Enable:                          types.BoolValue(service.ServerSideAdTracking.Enable),
			CheckAdMediaSegmentAvailability: types.BoolValue(service.ServerSideAdTracking.CheckAdMediaSegmentAvailability),
		}
	}

	// Add TranscodingProfile if provided
	if service.TranscodingProfile.Id != 0 {
		result.TranscodingProfile = &transcodingProfileDataSourceModel{
			ID:         types.Int64Value(int64(service.TranscodingProfile.Id)),
			Name:       types.StringValue(service.TranscodingProfile.Name),
			InternalId: types.StringValue(service.TranscodingProfile.InternalId),
			Content:    types.StringValue(service.TranscodingProfile.Content),
		}
	}

	// Add Source if provided
	if service.Source.Id != 0 {
		result.Source = &sourceLiteModel{
			ID:          types.Int64Value(int64(service.Source.Id)),
			Name:        types.StringValue(service.Source.Name),
			Type:        types.StringValue(service.Source.Type),
			Description: types.StringValue(service.Source.Description),
			MultiPeriod: types.BoolValue(service.Source.MultiPeriod),
			URL:         types.StringValue(service.Source.Url),
		}
	}

	// Add LiveAdReplacement if provided
	if service.LiveAdReplacement.AdServer.Id != 0 {
		result.LiveAdReplacement = &liveAdReplacementLiteModel{
			AdServer: adServerLiteModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.AdServer.Id)),
				Name: types.StringValue(service.LiveAdReplacement.AdServer.Name),
				Type: types.StringValue(service.LiveAdReplacement.AdServer.Type),
				URL:  types.StringValue(service.LiveAdReplacement.AdServer.Url),
			},
		}

		if service.LiveAdReplacement.GapFiller.Id != 0 {
			result.LiveAdReplacement.GapFiller = gapFillerModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.GapFiller.Id)),
				Name: types.StringValue(service.LiveAdReplacement.GapFiller.Name),
				Type: types.StringValue(service.LiveAdReplacement.GapFiller.Type),
				URL:  types.StringValue(service.LiveAdReplacement.GapFiller.Url),
			}
		}

		if service.LiveAdReplacement.SpotAware.Mode != "" {
			result.LiveAdReplacement.SpotAware = spotAwareModel{
				Mode: types.StringValue(service.LiveAdReplacement.SpotAware.Mode),
			}
		}
	}

	// Add LiveAdPreRoll if provided
	if service.LiveAdPreRoll.AdServer.Id != 0 {
		result.LiveAdPreRoll = &liveAdPrerollLiteModel{
			AdServer: adServerLiteModel{
				ID:   types.Int64Value(int64(service.LiveAdPreRoll.AdServer.Id)),
				Name: types.StringValue(service.LiveAdPreRoll.AdServer.Name),
				Type: types.StringValue(service.LiveAdPreRoll.AdServer.Type),
				URL:  types.StringValue(service.LiveAdPreRoll.AdServer.Url),
			},
			MaxDuration: types.Int64Value(int64(service.LiveAdPreRoll.MaxDuration)),
			Offset:      types.Int64Value(int64(service.LiveAdPreRoll.Offset)),
		}
	}

	if service.AdvancedOptions.AuthorizationHeader.Name != "" && service.AdvancedOptions.AuthorizationHeader.Value != "" {
		result.AdvancedOptions = &advancedOptionsModel{
			AuthorizationHeader: &authorizationHeaderModel{
				Name:  types.StringValue(service.AdvancedOptions.AuthorizationHeader.Name),
				Value: types.StringValue(service.AdvancedOptions.AuthorizationHeader.Value),
			},
		}
	}

	tflog.Debug(ctx, "Setting STATE ", map[string]interface{}{"state": result})

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serviceAdInsertionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state serviceAdInsertionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed adserver value from HashiCups
	service, err := r.client.GetAdInsertion(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", state.ID.ValueInt64(), err.Error()),
		)
		return
	}

	// Convert the []string to types.List
	tagsList, diags := types.ListValueFrom(ctx, types.StringType, service.Tags)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	state = serviceAdInsertionResourceModel{
		ID:                  types.Int64Value(int64(service.Id)),
		Name:                types.StringValue(service.Name),
		Type:                types.StringValue(service.Type),
		URL:                 types.StringValue(service.Url),
		CreationDate:        types.StringValue(service.CreationDate), // Make sure these fields exist in your API response
		UpdateDate:          types.StringValue(service.UpdateDate),
		State:               types.StringValue(service.State),
		Tags:                tagsList,
		EnableAdTranscoding: types.BoolValue(service.EnableAdTranscoding),
		ServerSideAdTracking: &serverSideAdTrackingModel{
			Enable:                          types.BoolValue(service.ServerSideAdTracking.Enable),
			CheckAdMediaSegmentAvailability: types.BoolValue(service.ServerSideAdTracking.CheckAdMediaSegmentAvailability),
		},
		Source: &sourceLiteModel{
			ID:          types.Int64Value(int64(service.Source.Id)),
			Name:        types.StringValue(service.Source.Name),
			Type:        types.StringValue(service.Source.Type),
			Description: types.StringValue(service.Source.Description),
			MultiPeriod: types.BoolValue(service.Source.MultiPeriod),
			URL:         types.StringValue(service.Source.Url),
		},
		TranscodingProfile: &transcodingProfileDataSourceModel{
			ID:         types.Int64Value(int64(service.TranscodingProfile.Id)),
			Name:       types.StringValue(service.TranscodingProfile.Name),
			InternalId: types.StringValue(service.TranscodingProfile.InternalId),
			Content:    types.StringValue(service.TranscodingProfile.Content),
		},
	}

	if service.LiveAdReplacement.AdServer.Id != 0 {
		state.LiveAdReplacement = &liveAdReplacementLiteModel{
			AdServer: adServerLiteModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.AdServer.Id)),
				Name: types.StringValue(service.LiveAdReplacement.AdServer.Name),
				Type: types.StringValue(service.LiveAdReplacement.AdServer.Type),
				URL:  types.StringValue(service.LiveAdReplacement.AdServer.Url),
			},
			GapFiller: gapFillerModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.GapFiller.Id)),
				Name: types.StringValue(service.LiveAdReplacement.GapFiller.Name),
				Type: types.StringValue(service.LiveAdReplacement.GapFiller.Type),
				URL:  types.StringValue(service.LiveAdReplacement.GapFiller.Url),
			},
			SpotAware: spotAwareModel{
				Mode: types.StringValue(service.LiveAdReplacement.SpotAware.Mode),
			},
		}
	}

	if service.LiveAdPreRoll.AdServer.Id != 0 {
		state.LiveAdPreRoll = &liveAdPrerollLiteModel{
			AdServer: adServerLiteModel{
				ID:   types.Int64Value(int64(service.LiveAdPreRoll.AdServer.Id)),
				Name: types.StringValue(service.LiveAdPreRoll.AdServer.Name),
				Type: types.StringValue(service.LiveAdPreRoll.AdServer.Type),
				URL:  types.StringValue(service.LiveAdPreRoll.AdServer.Url),
			},
			MaxDuration: types.Int64Value(int64(service.LiveAdPreRoll.MaxDuration)),
			Offset:      types.Int64Value(int64(service.LiveAdPreRoll.Offset)),
		}
	}

	if service.AdvancedOptions.AuthorizationHeader.Name != "" && service.AdvancedOptions.AuthorizationHeader.Value != "" {
		state.AdvancedOptions = &advancedOptionsModel{
			AuthorizationHeader: &authorizationHeaderModel{
				Name:  types.StringValue(service.AdvancedOptions.AuthorizationHeader.Name),
				Value: types.StringValue(service.AdvancedOptions.AuthorizationHeader.Value),
			},
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serviceAdInsertionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan and current state
	var plan serviceAdInsertionResourceModel
	// Get planned changes
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert from Terraform model to API model
	var tags []string
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var serviceData = broadpeakio.UpdateAdInsertionInput{
		Name: plan.Name.ValueString(),
		Tags: tags,
	}

	// Add TranscodingProfile if provided
	if plan.TranscodingProfile != nil {
		serviceData.TranscodingProfile = &broadpeakio.Identifiable{
			Id: uint(plan.TranscodingProfile.ID.ValueInt64()),
		}
	}

	// Add Source if provided
	if plan.Source != nil {
		serviceData.Source = &broadpeakio.Identifiable{
			Id: uint(plan.Source.ID.ValueInt64()),
		}
	}

	// Add EnableAdTranscoding if provided
	if !plan.EnableAdTranscoding.IsNull() {
		serviceData.EnableAdTranscoding = plan.EnableAdTranscoding.ValueBool()
	}

	// Add LiveAdPreRoll if provided
	if plan.LiveAdPreRoll != nil {
		serviceData.LiveAdPreRoll = &broadpeakio.LiveAdPreRoll{
			MaxDuration: uint(plan.LiveAdPreRoll.MaxDuration.ValueInt64()),
			Offset:      uint(plan.LiveAdPreRoll.Offset.ValueInt64()),
		}

		// Add AdServer if provided within LiveAdPreRoll
		if !plan.LiveAdPreRoll.AdServer.ID.IsNull() {
			serviceData.LiveAdPreRoll.AdServer = &broadpeakio.Identifiable{
				Id: uint(plan.LiveAdPreRoll.AdServer.ID.ValueInt64()),
			}
		}
	}

	// Add LiveAdReplacement if provided
	if plan.LiveAdReplacement != nil {
		serviceData.LiveAdReplacement = &broadpeakio.LiveAdReplacement{}

		// Add AdServer if provided within LiveAdReplacement
		if !plan.LiveAdReplacement.AdServer.ID.IsNull() {
			serviceData.LiveAdReplacement.AdServer = &broadpeakio.Identifiable{
				Id: uint(plan.LiveAdReplacement.AdServer.ID.ValueInt64()),
			}
		}

		// Add GapFiller if provided
		if !plan.LiveAdReplacement.GapFiller.ID.IsNull() {
			serviceData.LiveAdReplacement.GapFiller = &broadpeakio.Identifiable{
				Id: uint(plan.LiveAdReplacement.GapFiller.ID.ValueInt64()),
			}
		}

		// Add SpotAware if provided
		if !plan.LiveAdReplacement.SpotAware.Mode.IsNull() {
			serviceData.LiveAdReplacement.SpotAware = broadpeakio.SpotAware{
				Mode: plan.LiveAdReplacement.SpotAware.Mode.ValueString(),
			}
		}
	}

	// Add ServerSideAdTracking if provided
	if plan.ServerSideAdTracking != nil {
		serviceData.ServerSideAdTracking = &broadpeakio.ServerSideAdTracking{
			Enable:                          plan.ServerSideAdTracking.Enable.ValueBool(),
			CheckAdMediaSegmentAvailability: plan.ServerSideAdTracking.CheckAdMediaSegmentAvailability.ValueBool(),
		}
	}

	if plan.AdvancedOptions != nil && plan.AdvancedOptions.AuthorizationHeader != nil {
		serviceData.AdvancedOptions = &broadpeakio.AdvancedOptions{
			AuthorizationHeader: broadpeakio.AuthorizationHeader{
				Name:  plan.AdvancedOptions.AuthorizationHeader.Name.ValueString(),
				Value: plan.AdvancedOptions.AuthorizationHeader.Value.ValueString(),
			},
		}
	}

	// Retrieve ID from plan/state
	adinsertionID := uint(plan.ID.ValueInt64())

	tflog.Debug(ctx, "Update - Update Doc sent to BPKIO", map[string]interface{}{"id": adinsertionID, "updates": serviceData})

	// Update existing adserver
	_, err := r.client.UpdateAdInsertion(adinsertionID, serviceData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating adserver",
			"Could not update adserver, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetAdInsertion
	service, err := r.client.GetAdInsertion(adinsertionID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading AdInsertion",
			fmt.Sprintf("Could not fetch adinsertion service ID %d: %s", adinsertionID, err.Error()),
		)
		return
	}

	// Convert the []string to types.List
	tagsList, diags := types.ListValueFrom(ctx, types.StringType, service.Tags)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Map response body to schema and populate Computed attribute values
	result := serviceAdInsertionResourceModel{
		ID:                  types.Int64Value(int64(service.Id)),
		Name:                types.StringValue(service.Name),
		Type:                types.StringValue(service.Type),
		URL:                 types.StringValue(service.Url),
		CreationDate:        types.StringValue(service.CreationDate), // Make sure these fields exist in your API response
		UpdateDate:          types.StringValue(service.UpdateDate),
		State:               types.StringValue(service.State),
		Tags:                tagsList,
		EnableAdTranscoding: types.BoolValue(service.EnableAdTranscoding),
		Source: &sourceLiteModel{
			ID:          types.Int64Value(int64(service.Source.Id)),
			Name:        types.StringValue(service.Source.Name),
			Type:        types.StringValue(service.Source.Type),
			Description: types.StringValue(service.Source.Description),
			MultiPeriod: types.BoolValue(service.Source.MultiPeriod),
			URL:         types.StringValue(service.Source.Url),
		},
		TranscodingProfile: &transcodingProfileDataSourceModel{
			ID:         types.Int64Value(int64(service.TranscodingProfile.Id)),
			Name:       types.StringValue(service.TranscodingProfile.Name),
			InternalId: types.StringValue(service.TranscodingProfile.InternalId),
			Content:    types.StringValue(service.TranscodingProfile.Content),
		},
	}

	if service.LiveAdPreRoll.AdServer.Id != 0 {
		result.LiveAdPreRoll = &liveAdPrerollLiteModel{
			AdServer: adServerLiteModel{
				ID:   types.Int64Value(int64(service.LiveAdPreRoll.AdServer.Id)),
				Name: types.StringValue(service.LiveAdPreRoll.AdServer.Name),
				Type: types.StringValue(service.LiveAdPreRoll.AdServer.Type),
				URL:  types.StringValue(service.LiveAdPreRoll.AdServer.Url),
			},
			MaxDuration: types.Int64Value(int64(service.LiveAdPreRoll.MaxDuration)),
			Offset:      types.Int64Value(int64(service.LiveAdPreRoll.Offset)),
		}
	}

	if service.LiveAdReplacement.AdServer.Id != 0 {
		result.LiveAdReplacement = &liveAdReplacementLiteModel{
			AdServer: adServerLiteModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.AdServer.Id)),
				Name: types.StringValue(service.LiveAdReplacement.AdServer.Name),
				Type: types.StringValue(service.LiveAdReplacement.AdServer.Type),
				URL:  types.StringValue(service.LiveAdReplacement.AdServer.Url),
			},
			GapFiller: gapFillerModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.GapFiller.Id)),
				Name: types.StringValue(service.LiveAdReplacement.GapFiller.Name),
				Type: types.StringValue(service.LiveAdReplacement.GapFiller.Type),
				URL:  types.StringValue(service.LiveAdReplacement.GapFiller.Url),
			},
			SpotAware: spotAwareModel{
				Mode: types.StringValue(service.LiveAdReplacement.SpotAware.Mode),
			},
		}
	}

	// Add ServerSideAdTracking if provided
	if service.ServerSideAdTracking.Enable {
		result.ServerSideAdTracking = &serverSideAdTrackingModel{
			Enable:                          types.BoolValue(service.ServerSideAdTracking.Enable),
			CheckAdMediaSegmentAvailability: types.BoolValue(service.ServerSideAdTracking.CheckAdMediaSegmentAvailability),
		}
	}

	if service.AdvancedOptions.AuthorizationHeader.Name != "" && service.AdvancedOptions.AuthorizationHeader.Value != "" {
		result.AdvancedOptions = &advancedOptionsModel{
			AuthorizationHeader: &authorizationHeaderModel{
				Name:  types.StringValue(service.AdvancedOptions.AuthorizationHeader.Name),
				Value: types.StringValue(service.AdvancedOptions.AuthorizationHeader.Value),
			},
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, result)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serviceAdInsertionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state serviceAdInsertionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing adserver
	_, err := r.client.DeleteAdInsertion(uint(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Source AdServer",
			"Could not delete adserver, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state from the ID.
func (r *serviceAdInsertionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// serviceModel maps service schema data.
type serviceAdInsertionResourceModel struct {
	ID                   types.Int64                        `tfsdk:"id"`
	Name                 types.String                       `tfsdk:"name"`
	Type                 types.String                       `tfsdk:"type"`
	URL                  types.String                       `tfsdk:"url"`
	CreationDate         types.String                       `tfsdk:"creation_date"`
	UpdateDate           types.String                       `tfsdk:"update_date"`
	State                types.String                       `tfsdk:"state"`
	Tags                 types.List                         `tfsdk:"tags"`
	AdvancedOptions      *advancedOptionsModel              `tfsdk:"advanced_options"`
	LiveAdPreRoll        *liveAdPrerollLiteModel            `tfsdk:"live_ad_preroll"`
	LiveAdReplacement    *liveAdReplacementLiteModel        `tfsdk:"live_ad_replacement"`
	EnableAdTranscoding  types.Bool                         `tfsdk:"enable_ad_transcoding"`
	ServerSideAdTracking *serverSideAdTrackingModel         `tfsdk:"server_side_ad_tracking"`
	Source               *sourceLiteModel                   `tfsdk:"source"`
	TranscodingProfile   *transcodingProfileDataSourceModel `tfsdk:"transcoding_profile"`
}

type sourceLiteModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	URL         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description"`
	Format      types.String `tfsdk:"format"`
	MultiPeriod types.Bool   `tfsdk:"multi_period"`
}

type liveAdPrerollLiteModel struct {
	AdServer    adServerLiteModel `tfsdk:"ad_server"`
	MaxDuration types.Int64       `tfsdk:"max_duration"`
	Offset      types.Int64       `tfsdk:"offset"`
}

type liveAdReplacementLiteModel struct {
	AdServer  adServerLiteModel `tfsdk:"ad_server"`
	GapFiller gapFillerModel    `tfsdk:"gap_filler"`
	SpotAware spotAwareModel    `tfsdk:"spot_aware"`
}

type adServerLiteModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	URL  types.String `tfsdk:"url"`
}
