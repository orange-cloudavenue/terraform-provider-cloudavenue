package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetResourceConfig() map[testsacc.ResourceName]func() *testsacc.ResourceConfig {
	return map[testsacc.ResourceName]func() *testsacc.ResourceConfig{
		// * Catalog
		CatalogResourceName:               testsacc.NewResourceConfig(NewCatalogResourceTest()),
		CatalogACLResourceName:            testsacc.NewResourceConfig(NewCatalogACLResourceTest()),
		CatalogVAppTemplateDataSourceName: testsacc.NewResourceConfig(NewCatalogVAppTemplateDataSourceTest()),

		// * VDC
		VDCResourceName:      testsacc.NewResourceConfig(NewVDCResourceTest()),
		VDCGroupResourceName: testsacc.NewResourceConfig(NewVDCGroupResourceTest()),

		// * VAPP
		VAppResourceName:           testsacc.NewResourceConfig(NewVAppResourceTest()),
		VAppOrgNetworkResourceName: testsacc.NewResourceConfig(NewVAppOrgNetworkResourceTest()),

		// * Network
		NetworkRoutedResourceName: testsacc.NewResourceConfig(NewNetworkRoutedResourceTest()),

		// * Edge Gateway
		EdgeGatewayResourceName:               testsacc.NewResourceConfig(NewEdgeGatewayResourceTest()),
		EdgeGatewayAppPortProfileResourceName: testsacc.NewResourceConfig(NewEdgeGatewayAppPortProfileResourceTest()),
		EdgeGatewayFirewallResourceName:       testsacc.NewResourceConfig(NewEdgeGatewayFirewallResourceTest()),
		EdgeGatewaySecurityGroupResourceName:  testsacc.NewResourceConfig(NewEdgeGatewaySecurityGroupResourceTest()),
		EdgeGatewayDhcpForwardingResourceName: testsacc.NewResourceConfig(NewEdgeGatewayDhcpForwardingResourceTest()),

		// * Backup
		BackupResourceName: testsacc.NewResourceConfig(NewBackupResourceTest()),

		// * VM
		VMResourceName: testsacc.NewResourceConfig(NewVMResourceTest()),

		// * S3
		S3BucketResourceName:                        testsacc.NewResourceConfig(NewS3BucketResourceTest()),
		S3BucketVersioningConfigurationResourceName: testsacc.NewResourceConfig(NewS3BucketVersioningConfigurationResourceTest()),
		S3BucketCorsConfigurationResourceName:       testsacc.NewResourceConfig(NewS3BucketCorsConfigurationResourceTest()),
		S3BucketLifecycleConfigurationResourceName:  testsacc.NewResourceConfig(NewS3BucketLifecycleConfigurationResourceTest()),
		S3BucketWebsiteConfigurationResourceName:    testsacc.NewResourceConfig(NewS3BucketWebsiteConfigurationResourceTest()),
		S3BucketACLResourceName:                     testsacc.NewResourceConfig(NewS3BucketACLResourceTest()),
		S3CredentialResourceName:                    testsacc.NewResourceConfig(NewS3CredentialResourceTest()),
		S3BucketPolicyResourceName:                  testsacc.NewResourceConfig(NewS3BucketPolicyResourceTest()),

		// * VCDA
		VCDAIPResourceName: testsacc.NewResourceConfig(NewVCDAIPResourceTest()),

		// * Public IP
		PublicIPResourceName: testsacc.NewResourceConfig(NewPublicIPResourceTest()),

		// * IAM
		IAMUserResourceName: testsacc.NewResourceConfig(NewIAMUserResourceTest()),
	}
}
