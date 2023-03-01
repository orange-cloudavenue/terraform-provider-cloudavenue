// Package vdc provides a Terraform datasource.
package vdc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

const vdcStatuses = "`SAVING`, `SAVED`, `CONFIGURING`, `REALIZED`, `REALIZATION_FAILED`," +
	" `DELETING`, `DELETE_FAILED`, `OBJECT_NOT_FOUND`, `UNCONFIGURED`."

var (
	_ datasource.DataSource              = &vdcGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &vdcGroupDataSource{}
)

func NewVDCGroupDataSource() datasource.DataSource {
	return &vdcGroupDataSource{}
}

type vdcGroupDataSource struct {
	client *client.CloudAvenue
}

type vdcGroupDataSourceModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Description                types.String `tfsdk:"description"`
	ErrorMessage               types.String `tfsdk:"error_message"`
	DFWEnabled                 types.Bool   `tfsdk:"dfw_enabled"`
	LocalEgress                types.Bool   `tfsdk:"local_egress"`
	NetworkPoolID              types.String `tfsdk:"network_pool_id"`
	NetworkPoolUniversalID     types.String `tfsdk:"network_pool_universal_id"`
	NetworkProviderType        types.String `tfsdk:"network_provider_type"`
	Status                     types.String `tfsdk:"status"`
	Type                       types.String `tfsdk:"type"`
	UniversalNetworkingEnabled types.Bool   `tfsdk:"universal_networking_enabled"`
	Vdcs                       types.List   `tfsdk:"vdcs"`
}

type vdcModel struct {
	FaultDomainTag       types.String `tfsdk:"fault_domain_tag"`
	NetworkProviderScope types.String `tfsdk:"network_provider_scope"`
	IsRemoteOrg          types.Bool   `tfsdk:"is_remote_org"`
	Status               types.String `tfsdk:"status"`
	SiteID               types.String `tfsdk:"site_name"`
	SiteName             types.String `tfsdk:"site_id"`
	Name                 types.String `tfsdk:"name"`
	ID                   types.String `tfsdk:"id"`
}

var vdcModelAttrTypes = map[string]attr.Type{
	"fault_domain_tag":       types.StringType,
	"network_provider_scope": types.StringType,
	"is_remote_org":          types.BoolType,
	"status":                 types.StringType,
	"site_id":                types.StringType,
	"site_name":              types.StringType,
	"name":                   types.StringType,
	"id":                     types.StringType,
}

func (d *vdcGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "group"
}

func (d *vdcGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provides a data source to read vDC group and reference in other resources.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vDC Group.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the vDC Group.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the vDC Group.",
				Computed:            true,
			},
			"error_message": schema.StringAttribute{
				MarkdownDescription: "More detailed error message when vDC group has error status.",
				Computed:            true,
			},
			"dfw_enabled": schema.BoolAttribute{
				MarkdownDescription: "Distributed firewall status.",
				Computed:            true,
			},
			"local_egress": schema.BoolAttribute{
				MarkdownDescription: "Status whether local egress is enabled for a universal router belonging to a universal vDC group.",
				Computed:            true,
			},
			"network_pool_id": schema.StringAttribute{
				MarkdownDescription: "ID of used network pool of the vDC Group.",
				Computed:            true,
			},
			"network_pool_universal_id": schema.StringAttribute{
				MarkdownDescription: "The network provider’s universal id that is backing the universal network pool.",
				Computed:            true,
			},
			"network_provider_type": schema.StringAttribute{
				MarkdownDescription: "Defines the networking provider backing the vDC Group.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the group can be in " + vdcStatuses,
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the vDC Group (e.g. `LOCAL`, `UNIVERSAL`).",
				Computed:            true,
			},
			"universal_networking_enabled": schema.BoolAttribute{
				MarkdownDescription: "True means that a vDC group router has been created.",
				Computed:            true,
			},
			"vdcs": schema.ListNestedAttribute{
				MarkdownDescription: "The list of organization vDCs that are participating in this group.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"fault_domain_tag": schema.StringAttribute{
							MarkdownDescription: "Represents the fault domain of a given organization vDC.",
							Computed:            true,
						},
						"network_provider_scope": schema.StringAttribute{
							MarkdownDescription: "Specifies the network provider scope of the vDC.",
							Computed:            true,
						},
						"is_remote_org": schema.BoolAttribute{
							MarkdownDescription: "Specifies whether the vDC is local to this VCD site.",
							Computed:            true,
						},
						"status": schema.StringAttribute{
							MarkdownDescription: "The status of the vDC can be in " + vdcStatuses,
							Computed:            true,
						},
						"site_name": schema.StringAttribute{
							MarkdownDescription: "Site name that vDC belongs.",
							Computed:            true,
						},
						"site_id": schema.StringAttribute{
							MarkdownDescription: "Site ID that vDC belongs.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "vDC name.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "vDC ID.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *vdcGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *vdcGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data vdcGroupDataSourceModel
		diag diag.Diagnostics
	)
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	adminOrg, err := d.client.Vmware.GetAdminOrgByNameOrId(d.client.GetOrg())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving Org", err.Error())
		return
	}

	// Get the vdcGroup.
	vdcGroup, err := adminOrg.GetVdcGroupByName(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vDC Group", err.Error())
		return
	}

	// populate the data source model with the vdcGroup data
	data = vdcGroupDataSourceModel{
		ID:                         types.StringValue(vdcGroup.VdcGroup.Id),
		Name:                       types.StringValue(vdcGroup.VdcGroup.Name),
		Description:                types.StringValue(vdcGroup.VdcGroup.Description),
		DFWEnabled:                 types.BoolValue(vdcGroup.VdcGroup.DfwEnabled),
		ErrorMessage:               types.StringValue(vdcGroup.VdcGroup.ErrorMessage),
		LocalEgress:                types.BoolValue(vdcGroup.VdcGroup.LocalEgress),
		NetworkPoolID:              types.StringValue(vdcGroup.VdcGroup.NetworkPoolId),
		NetworkPoolUniversalID:     types.StringValue(vdcGroup.VdcGroup.NetworkPoolUniversalId),
		NetworkProviderType:        types.StringValue(vdcGroup.VdcGroup.NetworkProviderType),
		Status:                     types.StringValue(vdcGroup.VdcGroup.Status),
		Type:                       types.StringValue(vdcGroup.VdcGroup.Type),
		UniversalNetworkingEnabled: types.BoolValue(vdcGroup.VdcGroup.UniversalNetworkingEnabled),
	}

	listVdcs := make([]vdcModel, 0)
	for _, vdc := range vdcGroup.VdcGroup.ParticipatingOrgVdcs {
		listVdcs = append(listVdcs, vdcModel{
			FaultDomainTag:       types.StringValue(vdc.FaultDomainTag),
			NetworkProviderScope: types.StringValue(vdc.NetworkProviderScope),
			IsRemoteOrg:          types.BoolValue(vdc.RemoteOrg),
			Status:               types.StringValue(vdc.Status),
			SiteID:               types.StringValue(vdc.SiteRef.ID),
			SiteName:             types.StringValue(vdc.SiteRef.Name),
			ID:                   types.StringValue(vdc.VdcRef.ID),
			Name:                 types.StringValue(vdc.VdcRef.Name),
		})
	}

	data.Vdcs, diag = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: vdcModelAttrTypes}, listVdcs)

	resp.Diagnostics.Append(diag...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
