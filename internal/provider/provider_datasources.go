package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/alb"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/catalog"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/iam"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/publicip"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/storage"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vm"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vrf"
)

// DataSources defines the data sources implemented in the provider.
func (p *cloudavenueProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// * ALB
		alb.NewAlbPoolDataSource,

		// * TIER0
		vrf.NewTier0VrfsDataSource,
		vrf.NewTier0VrfDataSource,

		// * PUBLICIP
		publicip.NewPublicIPDataSource,

		// * EDGE GATEWAY
		edgegw.NewEdgeGatewayDataSource,
		edgegw.NewEdgeGatewaysDataSource,
		edgegw.NewFirewallDataSource,
		edgegw.NewSecurityGroupDataSource,
		edgegw.NewIPSetDataSource,
		edgegw.NewDhcpForwardingDataSource,
		edgegw.NewStaticRouteDataSource,
		edgegw.NewNATRuleDataSource,

		// * VDC
		vdc.NewVDCsDataSource,
		vdc.NewVDCDataSource,
		vdc.NewVDCGroupDataSource,

		// * VAPP
		vapp.NewVappDataSource,
		vapp.NewOrgNetworkDataSource,
		vapp.NewIsolatedNetworkDataSource,

		// * CATALOG
		catalog.NewCatalogsDataSource,
		catalog.NewCatalogDataSource,
		catalog.NewVAppTemplateDataSource,
		catalog.NewCatalogMediaDataSource,
		catalog.NewCatalogMediasDataSource,

		// * IAM
		iam.NewUserDataSource,
		iam.NewRoleDataSource,
		iam.NewIAMRightDataSource,

		// * VM
		vm.NewVMAffinityRuleDatasource,
		vm.NewVMDataSource,

		// * NETWORK
		network.NewNetworkIsolatedDataSource,
		network.NewNetworkRoutedDataSource,
		network.NewDhcpDataSource,
		network.NewDhcpBindingDataSource,

		// * STORAGE
		storage.NewProfileDataSource,
		storage.NewProfilesDataSource,
	}
}
