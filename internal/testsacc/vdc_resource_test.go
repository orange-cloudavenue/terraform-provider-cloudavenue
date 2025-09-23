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
						resource.TestCheckResourceAttr(resourceName, "vcpu", "10"),
						resource.TestCheckResourceAttr(resourceName, "memory", "30"),

						// ! Deprecated fields - Maintain for backward compatibility
						resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
						resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
						resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Update - replace deprecated fields with new fields
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description" }}
							vcpu                  = 10
							memory                = 30
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
							resource.TestCheckResourceAttr(resourceName, "vcpu", "10"),
							resource.TestCheckResourceAttr(resourceName, "memory", "30"),

							// ! Deprecated fields - Maintain for backward compatibility
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "22000"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
						},
					},
					// Update description and add more VCPU
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ generate . "description" "longString"}}
							vcpu                  = 12
							memory                = 30
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
							resource.TestCheckResourceAttr(resourceName, "vcpu", "12"),
							resource.TestCheckResourceAttr(resourceName, "memory", "30"),

							// ! Deprecated fields - Maintain for backward compatibility
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "26400"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "30"),
						},
					},
					// Update disponibility_class
					// NOTE : invalid disponibility_class for current configuration
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ generate . "description" "longString"}}
							vcpu                  = 10
							memory                = 30
							billing_model         = "PAYG"
							disponibility_class   = "DUALROOM"
							service_class         = "STD"
							storage_billing_model = "PAYG"

							storage_profiles = [{
								class   = "gold"
								default = true
								limit   = 500
							  }]
						}`),
						TFAdvanced: testsacc.TFAdvanced{
							ExpectNonEmptyPlan: true,
							ExpectError:        regexp.MustCompile(`validation error`),
						},
					},
					// Update vcpu & memory
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description" }}
							vcpu                  = 12
							memory                = 50
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
							resource.TestCheckResourceAttr(resourceName, "vcpu", "12"),
							resource.TestCheckResourceAttr(resourceName, "memory", "50"),

							// ! Deprecated fields - Maintain for backward compatibility
							resource.TestCheckResourceAttr(resourceName, "cpu_speed_in_mhz", "2200"),
							resource.TestCheckResourceAttr(resourceName, "cpu_allocated", "26400"),
							resource.TestCheckResourceAttr(resourceName, "memory_allocated", "50"),
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
						vcpu                  = 12
						memory                = 50
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
						resource.TestCheckResourceAttr(resourceName, "vcpu", "12"),
						resource.TestCheckResourceAttr(resourceName, "memory", "50"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Update description and reduce VCPU
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_reserved" {
							name                  = {{ get . "name" }}
							description           = {{ generate . "description" "longString"}}
							vcpu                  = 10
							memory                = 50
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
							resource.TestCheckResourceAttr(resourceName, "vcpu", "10"),
							resource.TestCheckResourceAttr(resourceName, "memory", "50"),
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
						vcpu                  = 5
						memory                = 30
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
						resource.TestCheckResourceAttr(resourceName, "vcpu", "5"),
						resource.TestCheckResourceAttr(resourceName, "memory", "30"),
					},
				},
			}
		},
		// This is used to test vdc_group resource
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
						vcpu                  = 5
						memory                = 30
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
						resource.TestCheckResourceAttr(resourceName, "vcpu", "5"),
						resource.TestCheckResourceAttr(resourceName, "memory", "30"),
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
				},
				// ! Set 2 storages profiles
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdc" "example_storage_profiles" {
						name                  = {{ generate . "name" }}
						description           = {{ generate . "description" "longString"}}
						vcpu                  = 5
						memory                = 30
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
						resource.TestCheckResourceAttr(resourceName, "storage_profiles.#", "2"),
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
					// Remove one storage profile
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							vcpu                  = 5
							memory                = 30
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
							resource.TestCheckResourceAttr(resourceName, "storage_profiles.#", "1"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "storage_profiles.*", map[string]string{
								"class":   "gold",
								"limit":   "500",
								"default": "true",
							}),
						},
					},
					// Update storage profile class to custom class
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							vcpu                  = 5
							memory                = 30
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
							ExpectError: regexp.MustCompile(`failed to add storage`),
						},
					},
					// Update storage profile class limit under minimum size
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							vcpu                  = 5
							memory                = 30
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
							ExpectError: regexp.MustCompile(`validation error`),
						},
					},
					// Update storage profile class limit over maximum size
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							vcpu                  = 5
							memory                = 30
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
							ExpectError: regexp.MustCompile(`validation error`),
						},
					},
					// Update storage profile class to valid size
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdc" "example_storage_profiles" {
							name                  = {{ get . "name" }}
							description           = {{ get . "description"}}
							vcpu                  = 5
							memory                = 30
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
							resource.TestCheckResourceAttr(resourceName, "storage_profiles.#", "2"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "storage_profiles.*", map[string]string{
								"class":   "silver",
								"default": "false",
								"limit":   "500",
							}),
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
