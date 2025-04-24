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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &OrgResource{}

const (
	OrgResourceName = testsacc.ResourceName("cloudavenue_org")
)

type OrgResource struct{}

func NewOrgResourceTest() testsacc.TestACC {
	return &OrgResource{}
}

// GetResourceName returns the name of the resource.
func (r *OrgResource) GetResourceName() string {
	return OrgResourceName.String()
}

func (r *OrgResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return
}

func (r *OrgResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					import {
						to = cloudavenue_org.example
						id = "example"
					}

					resource "cloudavenue_org" "example" {
						description = {{ generate . "description" }}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrSet(resourceName, "internet_billing_mode"),
						// email and name are not set in the template, so they should be empty
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org" "example" {
							name = {{ generate . "name" }}
							description = {{ generate . "description" }}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttrSet(resourceName, "internet_billing_mode"),
							// email are not set in the template, so they should be empty
						},
					},
					// * This step resets the resource to its initial values
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org" "example" {
							name = "Provider Terraform"
							description = "Provider Terraform"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", "Provider Terraform"),
							resource.TestCheckResourceAttr(resourceName, "description", "Provider Terraform"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateID:     "example",
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
	}
}

func TestAccOrgResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&OrgResource{}),
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
	})
}
