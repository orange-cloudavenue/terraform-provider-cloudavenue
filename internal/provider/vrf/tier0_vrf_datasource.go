// Package tier0 provides a Terraform resource to manage Tier-0 VRFs.
package tier0

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

var (
	_ datasource.DataSource              = &tier0VrfDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfDataSource{}
)

// NewTier0VrfDataSource returns a new datasource implementing the tier0_vrf data source.
func NewTier0VrfDataSource() datasource.DataSource {
	return &tier0VrfDataSource{}
}

type tier0VrfDataSource struct {
	client *client.CloudAvenue
}

type tier0VrfDataSourceModel struct {
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Provider     types.String   `tfsdk:"tier0_provider"`
	ClassService types.String   `tfsdk:"class_service"`
	Services     []segmentModel `tfsdk:"services"`
}

type segmentModel struct {
	Service types.String `tfsdk:"service"`
	VLANID  types.String `tfsdk:"vlan_id"`
}

func (d *tier0VrfDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "vrf"
}

func (d *tier0VrfDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "The Tier-0 VRF",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "data source retrieve informations about a Tier-0 VRF.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Tier-0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"tier0_provider": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Tier-O provider info.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"class_service": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "List of tag of the Tier0 VRF.",
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"services": superschema.ListNestedAttribute{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "Services list of the Tier0-VRF.",
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"service": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Service of the segment.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"vlan_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "VLAN ID of the segment.",
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}.GetDataSource(ctx)
}

func (d *tier0VrfDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *tier0VrfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data tier0VrfDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tier0Detail, _, err := d.client.APIClient.Tier0Api.GetTier0VrfByName(d.client.Auth, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Tier0 detail, got error: %s", err))
		return
	}

	data.Provider = types.StringValue(tier0Detail.Tier0Provider)
	if tier0Detail.ClassService != nil {
		data.ClassService = types.StringValue(string(*tier0Detail.ClassService))
	}

	if tier0Detail.Services != nil {
		for _, segment := range *tier0Detail.Services {
			data.Services = append(data.Services, segmentModel{
				Service: types.StringValue(segment.Service),
				VLANID:  types.StringValue(segment.VlanId),
			})
		}
	}

	// Generate ID for the data source
	data.ID = data.Name

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
