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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &serviceAdInsertionDataSource{}
	_ datasource.DataSourceWithConfigure = &serviceAdInsertionDataSource{}
)

// serviceAdInsertionDataSource is the data source implementation.
type serviceAdInsertionDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewServiceAdInsertionDataSource is a helper function to simplify the provider implementation.
func NewServiceAdInsertionDataSource() datasource.DataSource {
	return &serviceAdInsertionDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *serviceAdInsertionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *serviceAdInsertionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_ad_insertion"
}

// Schema defines the schema for the data source.
func (d *serviceAdInsertionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"live_ad_preroll": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"ad_server": schema.SingleNestedAttribute{
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
						Computed: true,
					},
					"max_duration": schema.Int64Attribute{
						Computed: true,
						Optional: true,
					},
					"offset": schema.Int64Attribute{
						Computed: true,
						Optional: true,
					},
				},
				Computed: true,
			},
			"live_ad_replacement": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"ad_server": schema.SingleNestedAttribute{
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
						Computed: true,
					},
					"gap_filler": schema.SingleNestedAttribute{
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
						},
						Computed: true,
					},
					"spot_aware": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								Computed: true,
							},
						},
						Computed: true,
					},
				},
				Computed: true,
			},
			"advanced_options": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"authorization_header": schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed: true,
							},
							"value": schema.StringAttribute{
								Computed: true,
							},
						},
						Optional: true,
						Computed: true,
					},
				},
				Optional: true,
				Computed: true,
			},
			"enable_ad_transcoding": schema.BoolAttribute{
				Computed: true,
			},
			"server_side_ad_tracking": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable": schema.BoolAttribute{
						Computed: true,
					},
					"check_ad_media_segment_availability": schema.BoolAttribute{
						Computed: true,
					},
				},
				Computed: true,
			},
			"source": schema.SingleNestedAttribute{
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
					"format": schema.StringAttribute{
						Computed: true,
					},
					"description": schema.StringAttribute{
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
				Computed: true,
			},
			"transcoding_profile": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Computed: true,
					},
					"name": schema.StringAttribute{
						Computed: true,
					},
					"internal_id": schema.StringAttribute{
						Computed: true,
					},
					"content": schema.StringAttribute{
						Computed: true,
					},
				},
				Optional: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *serviceAdInsertionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state serviceAdInsertionDataSourceModel
	var serviceid int64

	diags := req.Config.GetAttribute(ctx, path.Root("id"), &serviceid)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Get the service from the API
	service, err := d.client.GetAdInsertion(uint(serviceid))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", serviceid, err.Error()),
		)
		return
	}

	// Convert the []string to types.List
	tagsList, diags := types.ListValueFrom(ctx, types.StringType, service.Tags)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	serviceState := serviceAdInsertionDataSourceModel{
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
		Source: &sourceModel{
			ID:          types.Int64Value(int64(service.Source.Id)),
			Name:        types.StringValue(service.Source.Name),
			Type:        types.StringValue(service.Source.Type),
			URL:         types.StringValue(service.Source.Url),
			Description: types.StringValue(service.Source.Description),
			MultiPeriod: types.BoolValue(service.Source.MultiPeriod),
			Origin: &originModel{
				CustomHeaders: func() []customHeadersModel {
					var headers []customHeadersModel
					for _, header := range service.Source.Origin.CustomHeaders {
						headers = append(headers, customHeadersModel{
							Name:  types.StringValue(header.Name),
							Value: types.StringValue(header.Value),
						})
					}
					return headers
				}(),
			},
		},
		TranscodingProfile: &transcodingProfileDataSourceModel{
			ID:         types.Int64Value(int64(service.TranscodingProfile.Id)),
			Name:       types.StringValue(service.TranscodingProfile.Name),
			InternalId: types.StringValue(service.TranscodingProfile.InternalId),
			Content:    types.StringValue(service.TranscodingProfile.Content),
		},
		AdvancedOptions: &advancedOptionsModel{
			AuthorizationHeader: &authorizationHeaderModel{
				Name:  types.StringValue(service.AdvancedOptions.AuthorizationHeader.Name),
				Value: types.StringValue(service.AdvancedOptions.AuthorizationHeader.Value),
			},
		},
	}

	if service.LiveAdReplacement.AdServer.Id != 0 {
		serviceState.LiveAdReplacement = &liveAdReplacementModel{
			AdServer: adServerModel{
				ID:   types.Int64Value(int64(service.LiveAdReplacement.AdServer.Id)),
				Name: types.StringValue(service.LiveAdReplacement.AdServer.Name),
				Type: types.StringValue(service.LiveAdReplacement.AdServer.Type),
				URL:  types.StringValue(service.LiveAdReplacement.AdServer.Url),
				QueryParameters: func() []queryParametersModel {
					var params []queryParametersModel
					for _, param := range service.LiveAdReplacement.AdServer.QueryParameters {
						params = append(params, queryParametersModel{
							Type:  types.StringValue(param.Type),
							Name:  types.StringValue(param.Name),
							Value: types.StringValue(param.Value),
						})
					}
					return params
				}(),
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
		serviceState.LiveAdPreRoll = &liveAdPrerollModel{
			AdServer: adServerModel{
				ID:   types.Int64Value(int64(service.LiveAdPreRoll.AdServer.Id)),
				Name: types.StringValue(service.LiveAdPreRoll.AdServer.Name),
				Type: types.StringValue(service.LiveAdPreRoll.AdServer.Type),
				URL:  types.StringValue(service.LiveAdPreRoll.AdServer.Url),
				QueryParameters: func() []queryParametersModel {
					var params []queryParametersModel
					for _, param := range service.LiveAdPreRoll.AdServer.QueryParameters {
						params = append(params, queryParametersModel{
							Type:  types.StringValue(param.Type),
							Name:  types.StringValue(param.Name),
							Value: types.StringValue(param.Value),
						})
					}
					return params
				}(),
			},
			MaxDuration: types.Int64Value(int64(service.LiveAdPreRoll.MaxDuration)),
			Offset:      types.Int64Value(int64(service.LiveAdPreRoll.Offset)),
		}
	}

	// Set state
	diags = resp.State.Set(ctx, &serviceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// serviceModel maps service schema data.
type serviceAdInsertionDataSourceModel struct {
	ID                   types.Int64                        `tfsdk:"id"`
	Name                 types.String                       `tfsdk:"name"`
	Type                 types.String                       `tfsdk:"type"`
	URL                  types.String                       `tfsdk:"url"`
	CreationDate         types.String                       `tfsdk:"creation_date"`
	UpdateDate           types.String                       `tfsdk:"update_date"`
	State                types.String                       `tfsdk:"state"`
	Tags                 types.List                         `tfsdk:"tags"`
	AdvancedOptions      *advancedOptionsModel              `tfsdk:"advanced_options"`
	LiveAdPreRoll        *liveAdPrerollModel                `tfsdk:"live_ad_preroll"`
	LiveAdReplacement    *liveAdReplacementModel            `tfsdk:"live_ad_replacement"`
	EnableAdTranscoding  types.Bool                         `tfsdk:"enable_ad_transcoding"`
	ServerSideAdTracking *serverSideAdTrackingModel         `tfsdk:"server_side_ad_tracking"`
	Source               *sourceModel                       `tfsdk:"source"`
	TranscodingProfile   *transcodingProfileDataSourceModel `tfsdk:"transcoding_profile"`
}

type serverSideAdTrackingModel struct {
	Enable                          types.Bool `tfsdk:"enable"`
	CheckAdMediaSegmentAvailability types.Bool `tfsdk:"check_ad_media_segment_availability"`
}

type sourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	URL         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description"`
	Format      types.String `tfsdk:"format"`
	MultiPeriod types.Bool   `tfsdk:"multi_period"`
	Origin      *originModel `tfsdk:"origin"`
}

type customHeadersModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type originModel struct {
	CustomHeaders []customHeadersModel `tfsdk:"custom_headers"`
}

type advancedOptionsModel struct {
	AuthorizationHeader *authorizationHeaderModel `tfsdk:"authorization_header"`
}

type authorizationHeaderModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type liveAdPrerollModel struct {
	AdServer    adServerModel `tfsdk:"ad_server"`
	MaxDuration types.Int64   `tfsdk:"max_duration"`
	Offset      types.Int64   `tfsdk:"offset"`
}

type liveAdReplacementModel struct {
	AdServer  adServerModel  `tfsdk:"ad_server"`
	GapFiller gapFillerModel `tfsdk:"gap_filler"`
	SpotAware spotAwareModel `tfsdk:"spot_aware"`
}

type spotAwareModel struct {
	Mode types.String `tfsdk:"mode"`
}

type gapFillerModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
	URL  types.String `tfsdk:"url"`
}

type adServerModel struct {
	ID              types.Int64            `tfsdk:"id"`
	Name            types.String           `tfsdk:"name"`
	Type            types.String           `tfsdk:"type"`
	URL             types.String           `tfsdk:"url"`
	QueryParameters []queryParametersModel `tfsdk:"query_parameters"`
}

type queryParametersModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
	Type  types.String `tfsdk:"type"`
}
