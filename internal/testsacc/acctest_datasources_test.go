package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetDataSourceConfig() map[testsacc.ResourceName]func() *testsacc.ResourceConfig {
	return map[testsacc.ResourceName]func() *testsacc.ResourceConfig{
		// * Catalog
		CatalogDataSourceName:             testsacc.NewResourceConfig(NewCatalogDataSourceTest()),
		CatalogACLDataSourceName:          testsacc.NewResourceConfig(NewCatalogACLDataSourceTest()),
		CatalogsDataSourceName:            testsacc.NewResourceConfig(NewCatalogsDataSourceTest()),
		CatalogVAppTemplateDataSourceName: testsacc.NewResourceConfig(NewCatalogVAppTemplateDataSourceTest()),

		// * Tier0
		Tier0VRFDataSourceName: testsacc.NewResourceConfig(NewTier0VRFDataSourceTest()),

		// * VDC
		VDCDataSourceName:                      testsacc.NewResourceConfig(NewVDCDataSourceTest()),
		VDCGroupDataSourceName:                 testsacc.NewResourceConfig(NewVDCGroupDataSourceTest()),
		VDCNetworkIsolatedDataSourceName:       testsacc.NewResourceConfig(NewVDCNetworkIsolatedDataSourceTest()),
		VDCGDynamicSecurityGroupDataSourceName: testsacc.NewResourceConfig(NewVDCGDynamicSecurityGroupDataSourceTest()),

		// * VDC Group
		VDCGDataSourceName:                testsacc.NewResourceConfig(NewVDCGDataSourceTest()),
		VDCGIPSetDataSourceName:           testsacc.NewResourceConfig(NewVDCGIPSetDataSourceTest()),
		VDCGNetworkIsolatedDataSourceName: testsacc.NewResourceConfig(NewVDCGNetworkIsolatedDataSourceTest()),
		VDCGSecurityGroupDataSourceName:   testsacc.NewResourceConfig(NewVDCGSecurityGroupDataSourceTest()),
		VDCGFirewallDataSourceName:        testsacc.NewResourceConfig(NewVDCGFirewallDataSourceTest()),

		// * Backup
		BackupDataSourceName: testsacc.NewResourceConfig(NewBackupDataSourceTest()),

		// * EdgeGateway
		EdgeGatewayDataSourceName:               testsacc.NewResourceConfig(NewEdgeGatewayDataSourceTest()),
		EdgeGatewaysDataSourceName:              testsacc.NewResourceConfig(NewEdgeGatewaysDataSourceTest()),
		EdgeGatewayFirewallDataSourceName:       testsacc.NewResourceConfig(NewEdgeGatewayFirewallDataSourceTest()),
		EdgeGatewayAppPortProfileDatasourceName: testsacc.NewResourceConfig(NewEdgeGatewayAppPortProfileDatasourceTest()),
		EdgeGatewaySecurityGroupDataSourceName:  testsacc.NewResourceConfig(NewEdgeGatewaySecurityGroupDataSourceTest()),
		EdgeGatewayDhcpForwardingDataSourceName: testsacc.NewResourceConfig(NewEdgeGatewayDhcpForwardingDataSourceTest()),
		EdgeGatewayNATRuleDataSourceName:        testsacc.NewResourceConfig(NewEdgeGatewayNATRuleDataSourceTest()),
		EdgeGatewayIPSetDataSourceName:          testsacc.NewResourceConfig(NewEdgeGatewayIPSetDataSourceTest()),

		// * S3
		S3BucketVersioningConfigurationDatasourceName: testsacc.NewResourceConfig(NewS3BucketVersioningConfigurationDatasourceTest()),
		S3BucketDatasourceName:                        testsacc.NewResourceConfig(NewS3BucketDatasourceTest()),
		S3BucketCorsConfigurationDataSourceName:       testsacc.NewResourceConfig(NewS3BucketCorsConfigurationDataSourceTest()),
		S3BucketLifecycleConfigurationDataSourceName:  testsacc.NewResourceConfig(NewS3BucketLifecycleConfigurationDataSourceTest()),
		S3BucketWebsiteConfigurationDataSourceName:    testsacc.NewResourceConfig(NewS3BucketWebsiteConfigurationDataSourceTest()),
		S3BucketACLDataSourceName:                     testsacc.NewResourceConfig(NewS3BucketACLDataSourceTest()),
		S3BucketPolicyDataSourceName:                  testsacc.NewResourceConfig(NewS3BucketPolicyDataSourceTest()),
		S3UserDataSourceName:                          testsacc.NewResourceConfig(NewS3UserDataSourceTest()),

		// * Public IP
		PublicIPsDataSourceName: testsacc.NewResourceConfig(NewPublicIPsDataSourceTest()),

		// * IAM
		IAMRolesDataSourceName: testsacc.NewResourceConfig(NewIAMRolesDataSourceTest()),
		IAMUserDataSourceName:  testsacc.NewResourceConfig(NewIAMUserDataSourceTest()),

		// * VApp
		VAppDatasourceName:                testsacc.NewResourceConfig(NewVAppDatasourceTest()),
		VAppIsolatedNetworkDataSourceName: testsacc.NewResourceConfig(NewVAppIsolatedNetworkDataSourceTest()),
	}
}
