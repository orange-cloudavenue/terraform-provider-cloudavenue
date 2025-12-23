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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCResource{}

const (
	VDCResourceName = testsacc.ResourceName("cloudavenue_vdc")
)

type VDCResource struct{}

func NewVDCResourceTest() testsacc.TestACC {
	return &VDCResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCResource) GetResourceName() string {
	return VDCResourceName.String()
}

func (r *VDCResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *VDCResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDC)),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc" "example" {
						name                  = {{ generate . "name" }}
						description           = {{ generate . "description" "longString"}}
						cpu_allocated         = 22000
						memory_allocated      = 30
						cpu_speed_in_mhz      = 2200
						billing_model         = "PAYG"
						disponibility_class   = "ONE-ROOM"
						service_class         = "STD"
						storage_billing_model = "PAYG"
					  
						storage_profiles = [{
						  class   = "gold"
						  default = true
						  limit   = 500
						}]
					  
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
						resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
						resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Update description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ generate . "description" "longString"}}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"
						  
							storage_profiles = [{
							  class   = "gold"
							  default = true
							  limit   = 500
							}]
						  
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
						},
					},
					// Update cpu_allocated
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description" }}
							cpu_allocated         = 22500
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 500
							  }]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22500"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
						},
					},
					// Update cpu_speed_in_mhz
					// NOTE : This generate resource replacement
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description" }}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2300
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"
						  
							storage_profiles = [{
							  class   = "gold"
							  default = true
							  limit   = 500
							}]
						  
						}`),
						TFAdvanced: testsacc.TFAdvanced{
							PlanOnly:           true,
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile(`CPU speed in MHz attribute is not valid`),
						},
					},
					// Update memory_allocated
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description" }}
							cpu_allocated         = 22500
							memory_allocated      = 40
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 500
							  }]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22500"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "40"),
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateID:     testsacc.GetValueFromTemplate(resourceName, "name"),
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
		"example_reserved": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDC)),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "RESERVED"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc" "example_reserved" {
						name                  = {{ generate . "name" }}
						description           = {{ generate . "description" "longString"}}
						cpu_allocated         = 22000
						memory_allocated      = 30
						cpu_speed_in_mhz      = 2200
						billing_model         = "RESERVED"
						disponibility_class   = "ONE-ROOM"
						service_class         = "STD"
						storage_billing_model = "PAYG"
					  
						storage_profiles = [{
						  class   = "gold"
						  default = true
						  limit   = 500
						}]
					  
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
						resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
						resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Update cpu_speed_in_mhz
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_reserved" {
							name                  = {{ get . "name" }}
							description           = {{ generate . "description" "longString"}}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2000
							billing_model         = "RESERVED"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"
						  
							storage_profiles = [{
							  class   = "gold"
							  default = true
							  limit   = 500
							}]
						  
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2000"),
						},
					},
				},
			}
		},
		// This is used to test vdc_group resource
		"example_vdc_group_1": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDC)),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc" "example_vdc_group_1" {
						name                  = {{ generate . "name" }}
						description           = {{ generate . "description" "longString"}}
						cpu_allocated         = 22000
						memory_allocated      = 30
						cpu_speed_in_mhz      = 2200
						billing_model         = "PAYG"
						disponibility_class   = "ONE-ROOM"
						service_class         = "STD"
						storage_billing_model = "PAYG"
					  
						storage_profiles = [{
						  class   = "gold"
						  default = true
						  limit   = 500
						}]
					  
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
						resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
						resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
					},
				},
			}
		},
		"example_vdc_group_2": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDC)),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.class", "gold"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.default", "true"),
					resource.TestCheckResourceAttr(resourceName, "storage_profiles.0.limit", "500"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc" "example_vdc_group_2" {
						name                  = {{ generate . "name" }}
						description           = {{ generate . "description" "longString"}}
						cpu_allocated         = 22000
						memory_allocated      = 30
						cpu_speed_in_mhz      = 2200
						billing_model         = "PAYG"
						disponibility_class   = "ONE-ROOM"
						service_class         = "STD"
						storage_billing_model = "PAYG"
					  
						storage_profiles = [{
						  class   = "gold"
						  default = true
						  limit   = 500
						}]
					  
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
						resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
						resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
					},
				},
			}
		},
		// Test storage profile
		"example_storage_profiles": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDC)),
					resource.TestCheckResourceAttr(resourceName, "billing_model", "PAYG"),
					resource.TestCheckResourceAttr(resourceName, "disponibility_class", "ONE-ROOM"),
					resource.TestCheckResourceAttr(resourceName, "service_class", "STD"),
					resource.TestCheckResourceAttr(resourceName, "storage_billing_model", "PAYG"),

					resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
					resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
					resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
				},
				// ! Set 2 storages profiles
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc" "example_storage_profiles" {
						name                  = {{ generate . "name" }}
						description           = {{ generate . "description" "longString"}}
						cpu_allocated         = 22000
						memory_allocated      = 30
						cpu_speed_in_mhz      = 2200
						billing_model         = "PAYG"
						disponibility_class   = "ONE-ROOM"
						service_class         = "STD"
						storage_billing_model = "PAYG"
					  
						storage_profiles = [
						{
							class   = "gold"
							default = true
							limit   = 500
						},
						{
							class   = "silver"
							default = false
							limit   = 500
						}
						]
					  
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "storage_profiles.*", map[string]string{
							"class":   "silver",
							"default": "false",
							"limit":   "500",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "storage_profiles.*", map[string]string{
							"class":   "gold",
							"limit":   "500",
							"default": "true",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// This generate a bug in API. Waiting for fix
					// {
					// 	TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					// 	resource "cloudavenue_vdc" "example_storage_profiles" {
					// 		name                  = {{ get . "name" }}
					// 		description           = {{ get . "description"}}
					// 		cpu_allocated         = 22000
					// 		memory_allocated      = 30
					// 		cpu_speed_in_mhz      = 2200
					// 		billing_model         = "PAYG"
					// 		disponibility_class   = "ONE-ROOM"
					// 		service_class         = "STD"
					// 		storage_billing_model = "PAYG"

					// 		storage_profiles = [{
					// 			class   = "gold"
					// 			default = true
					// 			limit   = 500
					// 		}]

					// 	}`),
					// 	Checks: []resource.TestCheckFunc{},
					// },
					// Update storage profile class to custom class
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 500
							},
							{
								class   = "silver_ocb0001234"
								default = false
								limit   = 500
							}]

						}`),
						Checks: []resource.TestCheckFunc{},
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							// ! Generate a error because the class is not valid for this organization
							ExpectError: regexp.MustCompile(`Error updating VDC`),
						},
					},
					// Update storage profile class limit under minimum size
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 80
							}]
						}`),
						Checks: []resource.TestCheckFunc{},
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							// ! Generate a error because the limit is under the minimum size
							ExpectError: regexp.MustCompile(`Storage profile limit attribute is not valid`),
						},
					},
					// Update storage profile class limit over maximum size
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 82000
							}]
						}`),
						Checks: []resource.TestCheckFunc{},
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							// ! Generate a error because the limit is over the maximum size
							ExpectError: regexp.MustCompile(`Storage profile limit attribute is not valid`),
						},
					},
					// Update storage profile class to valid size
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							cpu_allocated         = 22000
							memory_allocated      = 30
							cpu_speed_in_mhz      = 2200
							billing_model         = "PAYG"
							disponibility_class   = "ONE-ROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 800
							},
							{
								class   = "silver"
								default = false
								limit   = 500
							}]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "storage_profiles.*", map[string]string{
								"class":   "gold",
								"limit":   "800",
								"default": "true",
							}),
						},
					},
				},
			}
		},
	}
}

func TestAccVDCResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCResource{}),
	})
}
