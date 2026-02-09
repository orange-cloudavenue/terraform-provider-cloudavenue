/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
		VDCDataSourceName:                testsacc.NewResourceConfig(NewVDCDataSourceTest()),
		VDCNetworkIsolatedDataSourceName: testsacc.NewResourceConfig(NewVDCNetworkIsolatedDataSourceTest()),

		// * VDC Group
		VDCGDataSourceName:                     testsacc.NewResourceConfig(NewVDCGDataSourceTest()),
		VDCGIPSetDataSourceName:                testsacc.NewResourceConfig(NewVDCGIPSetDataSourceTest()),
		VDCGNetworkIsolatedDataSourceName:      testsacc.NewResourceConfig(NewVDCGNetworkIsolatedDataSourceTest()),
		VDCGSecurityGroupDataSourceName:        testsacc.NewResourceConfig(NewVDCGSecurityGroupDataSourceTest()),
		VDCGDynamicSecurityGroupDataSourceName: testsacc.NewResourceConfig(NewVDCGDynamicSecurityGroupDataSourceTest()),
		VDCGAppPortProfileDatasourceName:       testsacc.NewResourceConfig(NewVDCGAppPortProfileDatasourceTest()),
		VDCGFirewallDataSourceName:             testsacc.NewResourceConfig(NewVDCGFirewallDataSourceTest()),
		VDCGNetworkRoutedDataSourceName:        testsacc.NewResourceConfig(NewVDCGNetworkRoutedDataSourceTest()),

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
		EdgeGatewayNetworkRoutedDataSourceName:  testsacc.NewResourceConfig(NewEdgeGatewayNetworkRoutedDataSourceTest()),
		EdgeGatewayServicesDataSourceName:       testsacc.NewResourceConfig(NewEdgeGatewayServicesDataSourceTest()),
		EdgeGatewayVPNIPSecDataSourceName:       testsacc.NewResourceConfig(NewEdgeGatewayVPNIPSecDataSourceTest()),
		EdgeGatewayStaticRouteDataSourceName:    testsacc.NewResourceConfig(NewEdgeGatewayStaticRouteDataSourceTest()),

		// * EdgeGateway LoadBalancer (elb)
		ELBServiceEngineGroupDataSourceName:   testsacc.NewResourceConfig(NewELBServiceEngineGroupDataSourceTest()),
		ELBServiceEngineGroupsDataSourceName:  testsacc.NewResourceConfig(NewELBServiceEngineGroupsDataSourceTest()),
		ELBPoolDataSourceName:                 testsacc.NewResourceConfig(NewELBPoolDataSourceTest()),
		ELBVirtualServiceDataSourceName:       testsacc.NewResourceConfig(NewELBVirtualServiceDataSourceTest()),
		ELBPoliciesHTTPRequestDataSourceName:  testsacc.NewResourceConfig(NewELBPoliciesHTTPRequestDataSourceTest()),
		ELBPoliciesHTTPResponseDataSourceName: testsacc.NewResourceConfig(NewELBPoliciesHTTPResponseDataSourceTest()),
		ELBPoliciesHTTPSecurityDataSourceName: testsacc.NewResourceConfig(NewELBPoliciesHTTPSecurityDataSourceTest()),

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

		// * Org
		OrgCertificateLibraryDatasourceName: testsacc.NewResourceConfig(NewOrgCertificateLibraryDatasourceTest()),
	}
}
