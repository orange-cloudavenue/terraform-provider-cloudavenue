package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetResourceConfig() map[testsacc.ResourceName]func() resourceConfig {
	return map[testsacc.ResourceName]func() resourceConfig{
		// * Catalog
		CatalogResourceName:               NewResourceConfig(NewCatalogResourceTest()),
		CatalogACLResourceName:            NewResourceConfig(NewCatalogACLResourceTest()),
		CatalogVAppTemplateDataSourceName: NewResourceConfig(NewCatalogVAppTemplateDataSourceTest()),

		// * VDC
		VDCResourceName:      NewResourceConfig(NewVDCResourceTest()),
		VDCGroupResourceName: NewResourceConfig(NewVDCGroupResourceTest()),

		// * VAPP
		VAppResourceName:           NewResourceConfig(NewVAppResourceTest()),
		VAppOrgNetworkResourceName: NewResourceConfig(NewVAppOrgNetworkResourceTest()),

		// * Network
		NetworkRoutedResourceName: NewResourceConfig(NewNetworkRoutedResourceTest()),

		// * Edge Gateway
		EdgeGatewayResourceName:               NewResourceConfig(NewEdgeGatewayResourceTest()),
		EdgeGatewayAppPortProfileResourceName: NewResourceConfig(NewEdgeGatewayAppPortProfileResourceTest()),
		EdgeGatewayFirewallResourceName:       NewResourceConfig(NewEdgeGatewayFirewallResourceTest()),
		EdgeGatewaySecurityGroupResourceName:  NewResourceConfig(NewEdgeGatewaySecurityGroupResourceTest()),

		// * Backup
		BackupResourceName: NewResourceConfig(NewBackupResourceTest()),

		// * VM
		VMResourceName: NewResourceConfig(NewVMResourceTest()),

		// * S3
		S3BucketResourceName:                        NewResourceConfig(NewS3BucketResourceTest()),
		S3BucketVersioningConfigurationResourceName: NewResourceConfig(NewS3BucketVersioningConfigurationResourceTest()),
		S3BucketCorsConfigurationResourceName:       NewResourceConfig(NewS3BucketCorsConfigurationResourceTest()),
		S3BucketLifecycleConfigurationResourceName:  NewResourceConfig(NewS3BucketLifecycleConfigurationResourceTest()),
		S3BucketWebsiteConfigurationResourceName:    NewResourceConfig(NewS3BucketWebsiteConfigurationResourceTest()),
		S3BucketACLResourceName:                     NewResourceConfig(NewS3BucketACLResourceTest()),
		S3CredentialResourceName:                    NewResourceConfig(NewS3CredentialResourceTest()),
		S3BucketPolicyResourceName:                  NewResourceConfig(NewS3BucketPolicyResourceTest()),

		// * VCDA
		VCDAIPResourceName: NewResourceConfig(NewVCDAIPResourceTest()),

		// * Public IP
		PublicIPResourceName: NewResourceConfig(NewPublicIPResourceTest()),
	}
}
