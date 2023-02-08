package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	_ "github.com/orange-cloudavenue/cloudavenue-sdk-go"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

var (
	_ datasource.DataSource              = &tier0VrfDataSource{}
	_ datasource.DataSourceWithConfigure = &tier0VrfDataSource{}
)

func NewTier0VrfDataSource() datasource.DataSource {
	return &tier0VrfDataSource{}
}

type tier0VrfDataSource struct {
	client *CloudAvenueClient
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
	VlanID  types.String `tfsdk:"vlan_id"`
}

func (d *tier0VrfDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tier0_vrf"
}

func (d *tier0VrfDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve information about a Tier-0 VRF.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Tier-0 VRF.",
				Required:            true,
			},
			"tier0_provider": schema.StringAttribute{
				MarkdownDescription: "Tier-O provider info.",
				Computed:            true,
			},
			"class_service": schema.StringAttribute{
				MarkdownDescription: "List of tag of the Tier0 VRF.",
				Computed:            true,
			},
			"services": schema.ListNestedAttribute{
				MarkdownDescription: "Services list of the Tier0-VRF.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service": schema.StringAttribute{
							MarkdownDescription: "Service of the segment.",
							Computed:            true,
						},
						"vlan_id": schema.StringAttribute{
							MarkdownDescription: "VLAN ID of the segment.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *tier0VrfDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*CloudAvenueClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *CloudAvenueClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
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

	tier0Detail, _, err := d.client.Tier0Api.ApiCustomersV20Tier0VrfsTier0NameGet(d.client.auth, data.Name.ValueString())
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
				VlanID:  types.StringValue(segment.VlanId),
			})
		}
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Generate ID for the data source
	data.ID = utils.GenerateUUID(data.Name.String())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
