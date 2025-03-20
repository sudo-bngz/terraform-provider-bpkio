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
	_ datasource.DataSource              = &transcodingProfileDataSource{}
	_ datasource.DataSourceWithConfigure = &transcodingProfileDataSource{}
)

// transcodingProfileDataSource is the data source implementation.
type transcodingProfileDataSource struct {
	client *broadpeakio.BroadpeakClient
}

// NewTranscodingProfileDataSource is a helper function to simplify the provider implementation.
func NewTranscodingProfileDataSource() datasource.DataSource {
	return &transcodingProfileDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *transcodingProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *transcodingProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_transcoding_profile"
}

// Schema defines the schema for the data source.
func (d *transcodingProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
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
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *transcodingProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state transcodingProfileDataSourceModel
	var sourceid int64

	diags := req.Config.GetAttribute(ctx, path.Root("id"), &sourceid)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	// Get the source from the API
	profile, err := d.client.GetTranscodingProfile(uint(sourceid))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Single Service",
			fmt.Sprintf("Service with ID %d not found (%s)", sourceid, err.Error()),
		)
		return
	}

	profileState := transcodingProfileDataSourceModel{
		ID:         types.Int64Value(int64(profile.Id)),
		Name:       types.StringValue(profile.Name),
		Content:    types.StringValue(profile.Content),
		InternalId: types.StringValue(profile.InternalId),
	}

	// Set state
	diags = resp.State.Set(ctx, &profileState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// transcodingProfileDataSourceModel maps sources schema data.
type transcodingProfileDataSourceModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Content    types.String `tfsdk:"content"`
	InternalId types.String `tfsdk:"internal_id"`
}
