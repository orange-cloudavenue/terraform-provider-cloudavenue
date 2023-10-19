package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetResourceConfig() map[testsacc.ResourceName]func() resourceConfig {
	return map[testsacc.ResourceName]func() resourceConfig{
		// * Catalog
		CatalogResourceName:    NewResourceConfig(NewCatalogResourceTest()),
		CatalogACLResourceName: NewResourceConfig(NewCatalogACLResourceTest()),

		// * VDC
		VDCResourceName:      NewResourceConfig(NewVDCResourceTest()),
		VDCGroupResourceName: NewResourceConfig(NewVDCGroupResourceTest()),

		// * VAPP
		VAppResourceName:           NewResourceConfig(NewVAppResourceTest()),
		VAppOrgNetworkResourceName: NewResourceConfig(NewVAppOrgNetworkResourceTest()),

		// * Network
		NetworkRoutedResourceName: NewResourceConfig(NewNetworkRoutedResourceTest()),

		// * Edge Gateway
		EdgeGatewayResourceName: NewResourceConfig(NewEdgeGatewayResourceTest()),

		// * Backup
		BackupResourceName: NewResourceConfig(NewBackupResourceTest()),

		// * S3
		S3BucketResourceName:                        NewResourceConfig(NewS3BucketResourceTest()),
		S3BucketVersioningConfigurationResourceName: NewResourceConfig(NewS3BucketVersioningConfigurationResourceTest()),
	}
}
