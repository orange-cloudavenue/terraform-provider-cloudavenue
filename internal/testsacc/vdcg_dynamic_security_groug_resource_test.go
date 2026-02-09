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

var _ testsacc.TestACC = &VDCGDynamicSecurityGroupResource{}

const (
	VDCGDynamicSecurityGroupResourceName = testsacc.ResourceName("cloudavenue_vdcg_dynamic_security_group")
)

type VDCGDynamicSecurityGroupResource struct{}

func NewVDCGDynamicSecurityGroupResourceTest() testsacc.TestACC {
	return &VDCGDynamicSecurityGroupResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGDynamicSecurityGroupResource) GetResourceName() string {
	return VDCGDynamicSecurityGroupResourceName.String()
}

func (r *VDCGDynamicSecurityGroupResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	return resp
}

func (r *VDCGDynamicSecurityGroupResource) Tests(_ context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Create an empty dynamic security group
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return resp
				},
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
				},
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_dynamic_security_group" "example" {
						vdc_group_name  = cloudavenue_vdcg.example.name
						name            = {{ generate . "name" }}
						description     = {{ generate . "description" }}
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "criteria.#", "0"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Update name
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_dynamic_security_group" "example" {
							vdc_group_name  = cloudavenue_vdcg.example.name
							name            = {{ generate . "name" }}
							description     = {{ get . "description" }}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "criteria.#", "0"),
						},
					},
					// * Update description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg_dynamic_security_group" "example" {
							vdc_group_name  = cloudavenue_vdcg.example.name
							name            = {{ get . "name" }}
							description     = {{ generate . "description" }}
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "criteria.#", "0"),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
		"example_with_criteria": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonDependencies: func() (resp testsacc.DependenciesConfigResponse) {
					resp.Append(GetResourceConfig()[VDCGResourceName]().GetDefaultConfig)
					return resp
				},
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.SecurityGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vdc_group_id"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_dynamic_security_group" "example_with_criteria" {
						vdc_group_id  = cloudavenue_vdcg.example.id
						name            = {{ generate . "name" }}
						description     = {{ generate . "description" }}
						criteria = [
							{ # OR
								rules = [ 
									{ # AND
										type = "VM_NAME"
										value = "test"
										operator = "STARTS_WITH"
									},
									{ # AND
									 	type = "VM_NAME"
										value = "front"
										operator = "CONTAINS"
									}
								]
							},
							{ # OR
							 	rules = [ 
									{ # AND
										type = "VM_TAG"
										value = "test"
										operator = "STARTS_WITH"
									},
									{ # AND
									 	type = "VM_TAG"
										value = "web-front"
										operator = "CONTAINS"
									}
								]
							}
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "criteria.#", "2"),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
							"type":     "VM_NAME",
							"value":    "test",
							"operator": "STARTS_WITH",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
							"type":     "VM_NAME",
							"value":    "front",
							"operator": "CONTAINS",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.1.rules.*", map[string]string{
							"type":     "VM_TAG",
							"value":    "test",
							"operator": "STARTS_WITH",
						}),
						resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.1.rules.*", map[string]string{
							"type":     "VM_TAG",
							"value":    "web-front",
							"operator": "CONTAINS",
						}),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// * Add a new criteria
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_dynamic_security_group" "example_with_criteria" {
						vdc_group_id  = cloudavenue_vdcg.example.id
						name            = {{ get . "name" }}
						description     = {{ get . "description" }}
						criteria = [
							{ # OR
								rules = [ 
									{ # AND
										type = "VM_NAME"
										value = "test"
										operator = "STARTS_WITH"
									},
									{ # AND
									 	type = "VM_NAME"
										value = "front"
										operator = "CONTAINS"
									},
									{ # AND
									 	type = "VM_TAG"
										value = "prod"
										operator = "ENDS_WITH"
									},
								]
							},
							{ # OR
							 	rules = [ 
									{ # AND
										type = "VM_TAG"
										value = "test"
										operator = "STARTS_WITH"
									},
									{ # AND
									 	type = "VM_TAG"
										value = "web-front"
										operator = "CONTAINS"
									}
								]
							},
							{ # OR
								rules = [
									{ # AND
										type = "VM_TAG"
										value = "prod"
										operator = "STARTS_WITH"
									},
									{ # AND
										type = "VM_TAG"
										value = "test-xx"
										operator = "EQUALS"
									}
								]
							}
						]
					}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "criteria.#", "3"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
								"type":     "VM_NAME",
								"value":    "test",
								"operator": "STARTS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
								"type":     "VM_NAME",
								"value":    "front",
								"operator": "CONTAINS",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "prod",
								"operator": "ENDS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.1.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "test",
								"operator": "STARTS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.1.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "web-front",
								"operator": "CONTAINS",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.2.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "prod",
								"operator": "STARTS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.2.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "test-xx",
								"operator": "EQUALS",
							}),
						},
					},
					// * Update name and description
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg_dynamic_security_group" "example_with_criteria" {
						vdc_group_id  = cloudavenue_vdcg.example.id
						name            = {{ generate . "name" }}
						description     = {{ generate . "description" }}
						criteria = [
							{ # OR
								rules = [ 
									{ # AND
										type = "VM_NAME"
										value = "test"
										operator = "STARTS_WITH"
									},
									{ # AND
									 	type = "VM_NAME"
										value = "front"
										operator = "CONTAINS"
									},
									{ # AND
									 	type = "VM_TAG"
										value = "prod"
										operator = "ENDS_WITH"
									},
								]
							},
							{ # OR
							 	rules = [ 
									{ # AND
										type = "VM_TAG"
										value = "test"
										operator = "STARTS_WITH"
									},
									{ # AND
									 	type = "VM_TAG"
										value = "web-front"
										operator = "CONTAINS"
									}
								]
							},
							{ # OR
								rules = [
									{ # AND
										type = "VM_TAG"
										value = "prod"
										operator = "STARTS_WITH"
									},
									{ # AND
										type = "VM_TAG"
										value = "test-xx"
										operator = "EQUALS"
									}
								]
							}
						]
					}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "criteria.#", "3"),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
								"type":     "VM_NAME",
								"value":    "test",
								"operator": "STARTS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
								"type":     "VM_NAME",
								"value":    "front",
								"operator": "CONTAINS",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.0.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "prod",
								"operator": "ENDS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.1.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "test",
								"operator": "STARTS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.1.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "web-front",
								"operator": "CONTAINS",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.2.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "prod",
								"operator": "STARTS_WITH",
							}),
							resource.TestCheckTypeSetElemNestedAttrs(resourceName, "criteria.2.rules.*", map[string]string{
								"type":     "VM_TAG",
								"value":    "test-xx",
								"operator": "EQUALS",
							}),
						},
					},
				},
				// ! Imports testing
				Imports: []testsacc.TFImport{
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "id"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_id", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
					{
						ImportStateIDBuilder: []string{"vdc_group_name", "name"},
						ImportState:          true,
						ImportStateVerify:    true,
					},
				},
			}
		},
	}
}

func TestAccVDCGDynamicSecurityGroupResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGDynamicSecurityGroupResource{}),
	})
}
