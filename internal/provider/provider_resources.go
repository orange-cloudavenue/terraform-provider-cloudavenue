package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/alb"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/catalog"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/iam"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/publicip"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vcda"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vm"
)

// Resources defines the resources implemented in the provider.
func (p *cloudavenueProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// * ALB
		alb.NewAlbPoolResource,

		// * EDGE GATEWAY
		edgegw.NewEdgeGatewayResource,
		edgegw.NewFirewallResource,
		edgegw.NewPortProfilesResource,
		edgegw.NewSecurityGroupResource,
		edgegw.NewIPSetResource,
		edgegw.NewDhcpForwardingResource,
		edgegw.NewStaticRouteResource,
		edgegw.NewNATRuleResource,

		// * VDC
		vdc.NewVDCResource,
		vdc.NewACLResource,

		// * VCDA
		vcda.NewVCDAIPResource,

		// * PUBLICIP
		publicip.NewPublicIPResource,

		// * VAPP
		vapp.NewVappResource,
		vapp.NewOrgNetworkResource,
		vapp.NewIsolatedNetworkResource,
		vapp.NewACLResource,

		// * CATALOG
		catalog.NewCatalogResource,

		// * IAM
		iam.NewIAMUserResource,
		iam.NewRoleResource,

		// * VM
		vm.NewDiskResource,
		vm.NewVMResource,
		vm.NewInsertedMediaResource,
		vm.NewVMAffinityRuleResource,
		vm.NewSecurityTagResource,

		// * NETWORK
		network.NewNetworkRoutedResource,
		network.NewNetworkIsolatedResource,
		network.NewDhcpBindingResource,
		network.NewDhcpResource,
	}
}
