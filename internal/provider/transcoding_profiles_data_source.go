package provider

import (
	"context"
	"fmt"

	broadpeakio "github.com/bashou/bpkio-go-sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &transcodingProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &transcodingProfileDataSource{}
)

// transcodingProfilesDataSource is the data source implementation.
type transcodingProfilesDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewTranscodingProfilesDataSource is a helper function to simplify the provider implementation.
func NewTranscodingProfilesDataSource() datasource.DataSource {
	return &transcodingProfilesDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *transcodingProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *transcodingProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_profiles"
}

// Schema defines the schema for the data source.
func (d *transcodingProfilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"profiles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"content": schema.StringAttribute{
							Computed: true,
						},
						"internal_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *transcodingProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state transcodingProfilesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	transcoding_profiles, err := d.client.GetAllTranscodingProfiles(0, 2000)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Transcoding Profiles",
			err.Error(),
		)
		return
	}

	state.Profiles = []transcodingProfileDataSourceModel{}
	for _, profile := range transcoding_profiles {
		state.Profiles = append(state.Profiles, transcodingProfileDataSourceModel{
			ID:         types.Int64Value(int64(profile.Id)),
			Name:       types.StringValue(profile.Name),
			Content:    types.StringValue(profile.Content),
			InternalId: types.StringValue(profile.InternalId),
		})
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// transcodingProfileDataSourceModel maps the data source schema data.
type transcodingProfilesDataSourceModel struct {
	Profiles []transcodingProfileDataSourceModel `tfsdk:"profiles"`
}
