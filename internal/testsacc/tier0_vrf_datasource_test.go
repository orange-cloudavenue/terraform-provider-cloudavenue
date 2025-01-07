/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package tier0 provides the acceptance tests for the provider.
package testsacc

import (
	"context"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &Tier0VRFDataSource{}

const (
	Tier0VRFDataSourceName = testsacc.ResourceName("data.cloudavenue_tier0_vrf")
)

type Tier0VRFDataSource struct{}

func NewTier0VRFDataSourceTest() testsacc.TestACC {
	return &Tier0VRFDataSource{}
}

// GetResourceName returns the name of the resource.
func (r *Tier0VRFDataSource) GetResourceName() string {
	return Tier0VRFDataSourceName.String()
}

func (r *Tier0VRFDataSource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *Tier0VRFDataSource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: `
					data "cloudavenue_tier0_vrf" "example" {
						name = "prvrf01eocb0006205allsp01"
					}`,
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttrSet(resourceName, "id"),
						resource.TestMatchResourceAttr(resourceName, "name", regexp.MustCompile("^prvrf[0-9]{2}eocb[0-9]{7}allsp[0-9]{2}")),
						resource.TestCheckResourceAttr(resourceName, "class_service", "VRF_STANDARD"),
						resource.TestCheckResourceAttr(resourceName, "tier0_provider", "pr01e02t0sp16"),
						resource.TestCheckResourceAttrSet(resourceName, "services.#"),
					},
				},
			}
		},
	}
}

func TestAccTier0VrfDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&Tier0VRFDataSource{}),
	})
}
