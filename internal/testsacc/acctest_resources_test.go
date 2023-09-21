package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetResourceConfig() map[testsacc.ResourceName]func() resourceConfig {
	return map[testsacc.ResourceName]func() resourceConfig{
		// * Catalog
		CatalogResourceName:    NewResourceConfig(NewCatalogResourceTest()),
		CatalogACLResourceName: NewResourceConfig(NewCatalogACLResourceTest()),

		// * VDC
		VDCResourceName: NewResourceConfig(NewVDCResourceTest()),

		// * VAPP
		VAppResourceName:           NewResourceConfig(NewVAppResourceTest()),
		VAppOrgNetworkResourceName: NewResourceConfig(NewVAppOrgNetworkResourceTest()),

		// * Network
		NetworkRoutedResourceName: NewResourceConfig(NewNetworkRoutedResourceTest()),

		// * Edge Gateway
		EdgeGatewayResourceName: NewResourceConfig(NewEdgeGatewayResourceTest()),
	}
}
