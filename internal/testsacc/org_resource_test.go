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

	"github.com/orange-cloudavenue/common-go/regex"
	"github.com/orange-cloudavenue/common-go/urn"

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
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Org)),
					resource.TestMatchResourceAttr(resourceName, "name", regex.OrganizationNameRegex()),
					resource.TestCheckResourceAttrSet(resourceName, "enabled"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.%"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.vdc"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.catalog"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.vapp"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.vm_running"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.user"),
					resource.TestCheckResourceAttrSet(resourceName, "resources.disk"),
				},
				// ! Import testing with data source
				// Import is tested with data source because create is not possible
				// and import with name or id works with both
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					import {
						to = cloudavenue_org.example
						id = "myOrganizationName"
					}

					resource "cloudavenue_org" "example" {
						description = {{ generate . "description" }}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttrSet(resourceName, "internet_billing_mode"),
						resource.TestCheckResourceAttrSet(resourceName, "full_name"),
						resource.TestCheckResourceAttrSet(resourceName, "email"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org" "example" {
							full_name = {{ generate . "fullname" }}
							description = {{ generate . "description" }}
							email = "foo@bar.com"
							internet_billing_mode = "TRAFFIC_VOLUME"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "full_name", testsacc.GetValueFromTemplate(resourceName, "fullname")),
							resource.TestCheckResourceAttr(resourceName, "email", "foo@bar.com"),
							resource.TestCheckResourceAttr(resourceName, "internet_billing_mode", "TRAFFIC_VOLUME"),
						},
					},
					// Second update to test multiple updates
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_org" "example" {
							full_name = "Mike and Dave Corporation"
							description = "Managed by terraform"
							email = "bar@foo.com"
							internet_billing_mode = "PAYG"
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", "Managed by terraform"),
							resource.TestCheckResourceAttr(resourceName, "full_name", "Mike and Dave Corporation"),
							resource.TestCheckResourceAttr(resourceName, "email", "bar@foo.com"),
							resource.TestCheckResourceAttr(resourceName, "internet_billing_mode", "PAYG"),
						},
					},
				},

				// ! Delete testing
				// Delete is not possible for an organization
				Destroy: false,

				// ! Imports testing
				// Import with both name and id
				// even if these attributes are informational here
				// And by simplicity we recommanded to use the name of the organization
				Imports: []testsacc.TFImport{
					{
						ImportStateID:     testsacc.GetValueFromTemplate(resourceName, "id"),
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						ImportStateID:     testsacc.GetValueFromTemplate(resourceName, "name"),
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
