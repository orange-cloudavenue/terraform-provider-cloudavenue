// Package vapp provides a Terraform datasource.
package vapp

import (
	"context"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

var (
	_ datasource.DataSource              = &isolatedNetworkDataSource{}
	_ datasource.DataSourceWithConfigure = &isolatedNetworkDataSource{}
)

func NewIsolatedNetworkDataSource() datasource.DataSource {
	return &isolatedNetworkDataSource{}
}

type isolatedNetworkDataSource struct {
	client *client.CloudAvenue

	org  org.Org
	vdc  vdc.VDC
	vapp vapp.VAPP
}

// Init Initializes the data source.
func (d *isolatedNetworkDataSource) Init(ctx context.Context, dm *isolatedNetworkDataSourceModel) (diags diag.Diagnostics) {
	d.org, diags = org.Init(d.client)
	if diags.HasError() {
		return
	}

	d.vdc, diags = vdc.Init(d.client, dm.VDC)
	if diags.HasError() {
		return
	}

	d.vapp, diags = vapp.Init(d.client, d.vdc, dm.VAppID, dm.VAppName)

	return
}

func (d *isolatedNetworkDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_" + "isolated_network"
}

func (d *isolatedNetworkDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = network.GetSchema(network.SetIsolatedVapp()).GetDataSource(ctx)
}

func (d *isolatedNetworkDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *isolatedNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_vapp_isolated_network", d.client.GetOrgName(), metrics.Read)()

	var (
		config      = &isolatedNetworkDataSourceModel{}
		diag        diag.Diagnostics
		vAppNetwork = govcdtypes.VAppNetworkConfiguration{}
	)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get vApp Network information
	vAppNetworkConfig, err := d.vapp.GetNetworkConfig()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vApp network config", err.Error())
		return
	}

	for _, networkConfig := range vAppNetworkConfig.NetworkConfig {
		if networkConfig.NetworkName == config.Name.ValueString() {
			vAppNetwork = networkConfig
		}
	}

	if vAppNetwork == (govcdtypes.VAppNetworkConfiguration{}) {
		resp.State.RemoveResource(ctx)
		return
	}

	// Get UUID.
	networkID, err := govcd.GetUuidFromHref(vAppNetwork.Link.HREF, false)
	if err != nil {
		resp.Diagnostics.AddError("Error on getting vApp network ID", err.Error())
		return
	}

	plan := &isolatedNetworkDataSourceModel{
		ID:                 utils.StringValueOrNull(uuid.Normalize(uuid.Network, networkID).String()),
		VDC:                utils.StringValueOrNull(d.vdc.GetName()),
		Name:               utils.StringValueOrNull(vAppNetwork.NetworkName),
		Description:        utils.StringValueOrNull(vAppNetwork.Description),
		VAppName:           utils.StringValueOrNull(config.VAppName.ValueString()),
		VAppID:             utils.StringValueOrNull(config.VAppID.ValueString()),
		Netmask:            types.StringNull(),
		Gateway:            types.StringNull(),
		DNS1:               types.StringNull(),
		DNS2:               types.StringNull(),
		DNSSuffix:          types.StringNull(),
		GuestVLANAllowed:   types.BoolValue(*vAppNetwork.Configuration.GuestVlanAllowed),
		RetainIPMacEnabled: types.BoolValue(*vAppNetwork.Configuration.RetainNetInfoAcrossDeployments),
	}

	// Get IP Scopes
	if len(vAppNetwork.Configuration.IPScopes.IPScope) > 0 {
		plan.Netmask = utils.StringValueOrNull(vAppNetwork.Configuration.IPScopes.IPScope[0].Netmask)
		plan.Gateway = utils.StringValueOrNull(vAppNetwork.Configuration.IPScopes.IPScope[0].Gateway)
		plan.DNS1 = utils.StringValueOrNull(vAppNetwork.Configuration.IPScopes.IPScope[0].DNS1)
		plan.DNS2 = utils.StringValueOrNull(vAppNetwork.Configuration.IPScopes.IPScope[0].DNS2)
		plan.DNSSuffix = utils.StringValueOrNull(vAppNetwork.Configuration.IPScopes.IPScope[0].DNSSuffix)
	}

	// Loop on static_ip_pool if it is not nil
	staticIPRanges := make([]staticIPPoolModel, 0)
	plan.StaticIPPool = types.SetNull(types.ObjectType{AttrTypes: staticIPPoolModelAttrTypes})
	if vAppNetwork.Configuration.IPScopes.IPScope[0].IPRanges != nil {
		for _, staticIPRange := range vAppNetwork.Configuration.IPScopes.IPScope[0].IPRanges.IPRange {
			staticIPRanges = append(staticIPRanges, staticIPPoolModel{
				StartAddress: utils.StringValueOrNull(staticIPRange.StartAddress),
				EndAddress:   utils.StringValueOrNull(staticIPRange.EndAddress),
			})
		}

		plan.StaticIPPool, diag = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: staticIPPoolModelAttrTypes}, staticIPRanges)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
