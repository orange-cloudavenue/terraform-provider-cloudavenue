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

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VCDAIPResource{}

const (
	VCDAIPResourceName = testsacc.ResourceName("cloudavenue_vcda_ip")
)

type VCDAIPResource struct{}

func NewVCDAIPResourceTest() testsacc.TestACC {
	return &VCDAIPResource{}
}

// GetResourceName returns the name of the resource.
func (r *VCDAIPResource) GetResourceName() string {
	return VCDAIPResourceName.String()
}

func (r *VCDAIPResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *VCDAIPResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * First test named "example"
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VCDA)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vcda_ip" "example" {
						ip_address = "10.0.0.1"
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.0.1"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// No updates
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"ip_address"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
				Destroy: true,
			}
		},
		"example_multiple": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VCDA)),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vcda_ip" "example_multiple" {
						ip_address = "10.0.0.1"
					}
					
					resource "cloudavenue_vcda_ip" "example_multiple-2" {
						ip_address = "10.0.0.2"
					}

					resource "cloudavenue_vcda_ip" "example_multiple-3" {
						ip_address = "10.0.0.3"
					}
					
					`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.0.1"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// No updates
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"ip_address"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVCDAIPResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VCDAIPResource{}),
	})
}
