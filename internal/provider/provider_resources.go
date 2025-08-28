/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/backup"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/catalog"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/draas"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/elb"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/iam"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/network"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/org"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/publicip"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/s3"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vcda"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vdcg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/vm"
)

// Resources defines the resources implemented in the provider.
func (p *cloudavenueProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// * EdgeGateway
		edgegw.NewEdgeGatewayResource,
		edgegw.NewFirewallResource,
		edgegw.NewAppPortProfileResource,
		edgegw.NewSecurityGroupResource,
		edgegw.NewIPSetResource,
		edgegw.NewDhcpForwardingResource,
		edgegw.NewStaticRouteResource,
		edgegw.NewNATRuleResource,
		edgegw.NewVPNIPSecResource,
		edgegw.NewNetworkRoutedResource,
		edgegw.NewServicesResource,

		// * EdgeGateway LoadBalancer
		elb.NewPoolResource,
		elb.NewVirtualServiceResource,
		elb.NewPoliciesHTTPRequestResource,
		elb.NewPoliciesHTTPResponseResource,
		elb.NewPoliciesHTTPSecurityResource,

		// * VDC
		vdc.NewVDCResource,
		vdc.NewACLResource,
		vdc.NewNetworkIsolatedResource,

		// * VDC Group
		vdcg.NewVDCGResource,
		vdcg.NewIPSetResource,
		vdcg.NewNetworkIsolatedResource,
		vdcg.NewDynamicSecurityGroupResource,
		vdcg.NewSecurityGroupResource,
		vdcg.NewAppPortProfileResource,
		vdcg.NewFirewallResource,
		vdcg.NewNetworkRoutedResource,

		// * DRAAS
		draas.NewDraasIPResource,
		// ! VCDA - Deprecated
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
		catalog.NewACLResource,

		// * IAM
		iam.NewIAMUserResource,
		iam.NewUserSAMLResource,
		iam.NewRoleResource,
		iam.NewTokenResource,

		// * VM
		vm.NewDiskResource,
		vm.NewVMResource,
		vm.NewInsertedMediaResource,
		vm.NewVMAffinityRuleResource,
		vm.NewSecurityTagResource,

		// * NETWORK
		network.NewNetworkRoutedResource,
		network.NewDhcpBindingResource,
		network.NewDhcpResource,

		// * BACKUP
		backup.NewBackupResource,

		// * S3
		s3.NewBucketVersioningConfigurationResource,
		s3.NewBucketResource,
		s3.NewBucketCorsConfigurationResource,
		s3.NewBucketLifecycleConfigurationResource,
		s3.NewBucketWebsiteConfigurationResource,
		s3.NewBucketACLResource,
		s3.NewCredentialResource,
		s3.NewBucketPolicyResource,

		// * ORG
		org.NewOrgResource,
		org.NewCertificateLibraryResource,
	}
}
