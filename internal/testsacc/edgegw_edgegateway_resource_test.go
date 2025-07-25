/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// package testsacc provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/edgegw"
)

var _ testsacc.TestACC = &EdgeGatewayResource{}

const (
	EdgeGatewayResourceName = testsacc.ResourceName("cloudavenue_edgegateway")
)

type EdgeGatewayResource struct{}

func NewEdgeGatewayResourceTest() testsacc.TestACC {
	return &EdgeGatewayResource{}
}

// GetResourceName returns the name of the resource.
func (r *EdgeGatewayResource) GetResourceName() string {
	return EdgeGatewayResourceName.String()
}

func (r *EdgeGatewayResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetDataSourceConfig()[Tier0VRFDataSourceName]().GetDefaultConfig)
	return
}

func (r *EdgeGatewayResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),

					// Read-Only attributes
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					resource.TestMatchResourceAttr(resourceName, "tier0_vrf_name", regexp.MustCompile(regexpTier0VRFName)),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway" "example" {
						owner_name     = cloudavenue_vdc.example.name
						tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
						bandwidth      = 25
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "bandwidth"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Test one of range value allowed in bandwidth attribute
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example" {
							owner_name     = cloudavenue_vdc.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							bandwidth      = 20
						  }`),
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile(`Invalid Bandwidth value`),
						},
					},
					// Test overcommit bandwidth
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example" {
							owner_name     = cloudavenue_vdc.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							bandwidth      = 300
						  }`),
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile(`Overcommitting bandwidth`),
						},
					},
					// Update bandwidth
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example" {
							owner_name     = cloudavenue_vdc.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							bandwidth      = 5
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "bandwidth", "5"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
					},
				},
				Destroy: true,
			}
		},
		"example_with_vdc_group": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),

					// Read-Only attributes
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					resource.TestMatchResourceAttr(resourceName, "tier0_vrf_name", regexp.MustCompile(regexpTier0VRFName)),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway" "example_with_vdc_group" {
						owner_name     = cloudavenue_vdcg.example.name
						tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
						bandwidth      = 25
					  }`),
					Checks: []resource.TestCheckFunc{},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_edgegateway" "example_with_vdc_group" {
							owner_name     = cloudavenue_vdcg.example.name
							tier0_vrf_name = data.cloudavenue_tier0_vrf.example.name
							bandwidth      = 5
						  }`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "bandwidth", "5"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_without_t_0": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrSet(resourceName, "owner_name"),

					// Read-Only attributes
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Gateway)),
					resource.TestMatchResourceAttr(resourceName, "tier0_vrf_name", regexp.MustCompile(regexpTier0VRFName)),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				},
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCResourceName]().GetDefaultConfig)
					return
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_edgegateway" "example_without_t_0" {
						owner_name     = cloudavenue_vdc.example.name
						bandwidth      = 5
					  }`),
					Checks: []resource.TestCheckFunc{},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccEdgeGatewayResource(t *testing.T) {
	edgegw.ConfigEdgeGateway = func() edgegw.EdgeGatewayConfig {
		return edgegw.EdgeGatewayConfig{
			CheckJobDelay: 10 * time.Second,
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&EdgeGatewayResource{}),
	})
}
