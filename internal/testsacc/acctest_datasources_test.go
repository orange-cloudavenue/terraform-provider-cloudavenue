package testsacc

import "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"

func GetDataSourceConfig() map[testsacc.ResourceName]func() resourceConfig {
	return map[testsacc.ResourceName]func() resourceConfig{
		Tier0VRFDataSourceName: NewResourceConfig(NewTier0VRFDataSourceTest()),

		// * VDC
		VDCDataSourceName:      NewResourceConfig(NewVDCDataSourceTest()),
		VDCGroupDataSourceName: NewResourceConfig(NewVDCGroupDataSourceTest()),

		// * Backup
		BackupDataSourceName: NewResourceConfig(NewBackupDataSourceTest()),

		// * EdgeGateway
		EdgeGatewayDataSourceName: NewResourceConfig(NewEdgeGatewayDataSourceTest()),

		// * S3
		S3BucketVersioningConfigurationDatasourceName: NewResourceConfig(NewS3BucketVersioningConfigurationDatasourceTest()),
		S3BucketDatasourceName:                        NewResourceConfig(NewS3BucketDatasourceTest()),
		S3BucketCorsConfigurationDataSourceName:       NewResourceConfig(NewS3BucketCorsConfigurationDataSourceTest()),
		S3BucketLifecycleConfigurationDataSourceName:  NewResourceConfig(NewS3BucketLifecycleConfigurationDataSourceTest()),
		S3BucketWebsiteConfigurationDataSourceName:    NewResourceConfig(NewS3BucketWebsiteConfigurationDataSourceTest()),
		S3BucketACLDataSourceName:                     NewResourceConfig(NewS3BucketACLDataSourceTest()),
		S3BucketPolicyDataSourceName:                  NewResourceConfig(NewS3BucketPolicyDataSourceTest()),
	}
}
