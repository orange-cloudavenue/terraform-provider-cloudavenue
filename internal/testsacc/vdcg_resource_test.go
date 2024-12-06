package testsacc

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/helpers/testsacc"
)

var _ testsacc.TestACC = &VDCGResource{}

const (
	VDCGResourceName = testsacc.ResourceName("cloudavenue_vdcg")
)

type VDCGResource struct{}

func NewVDCGResourceTest() testsacc.TestACC {
	return &VDCGResource{}
}

// GetResourceName returns the name of the resource.
func (r *VDCGResource) GetResourceName() string {
	return VDCGResourceName.String()
}

func (r *VDCGResource) DependenciesConfig() (resp testsacc.DependenciesConfigResponse) {
	resp.Append(GetResourceConfig()[VDCResourceName]().GetSpecificConfig("example_vdc_group_1"))
	resp.Append(GetResourceConfig()[VDCResourceName]().GetSpecificConfig("example_vdc_group_2"))
	return
}

func (r *VDCGResource) Tests(ctx context.Context) map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test {
	return map[testsacc.TestName]func(ctx context.Context, resourceName string) testsacc.Test{
		// * Test One (example)
		"example": func(_ context.Context, resourceName string) testsacc.Test {
			return testsacc.Test{
				CommonChecks: []resource.TestCheckFunc{
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
				},
				// ! Create testing
				Create: testsacc.TFConfig{
					TFConfig: testsacc.GenerateFromTemplate(resourceName, `
					resource "cloudavenue_vdcg" "example" {
						name = {{ generate . "name" }}
						description = {{ generate . "description" }}
						vdc_ids = [
							cloudavenue_vdc.example_vdc_group_1.id,
						]
					}`),
					Checks: []resource.TestCheckFunc{
						resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
						resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
						resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "1"),
					},
				},
				// ! Updates testing
				Updates: []testsacc.TFConfig{
					// Update description and add a new vdc_id
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg" "example" {
							name = {{ get . "name" }}
							description = {{ generate . "description" }}
							vdc_ids = [
								cloudavenue_vdc.example_vdc_group_1.id,
								cloudavenue_vdc.example_vdc_group_2.id,
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "2"),
						},
					},
					{
						TFConfig: testsacc.GenerateFromTemplate(resourceName, `
						resource "cloudavenue_vdcg" "example" {
							name = {{ generate . "name" }}
							description = {{ get . "description" }}
							vdc_ids = [
								cloudavenue_vdc.example_vdc_group_1.id,
								cloudavenue_vdc.example_vdc_group_2.id,
							]
						}`),
						Checks: []resource.TestCheckFunc{
							resource.TestCheckResourceAttr(resourceName, "description", testsacc.GetValueFromTemplate(resourceName, "description")),
							resource.TestCheckResourceAttr(resourceName, "name", testsacc.GetValueFromTemplate(resourceName, "name")),
							resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "2"),
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
					{
						ImportStateID:     testsacc.GetValueFromTemplate(resourceName, "id"),
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			}
		},
		// This test not working on testsacc
		// "example_move_state": func(_ context.Context, resourceName string) testsacc.Test {
		// 	return testsacc.Test{
		// 		CommonChecks: []resource.TestCheckFunc{},
		// 		// ! Create testing
		// 		Create: testsacc.TFConfig{
		// 			TFConfig: testsacc.GenerateFromTemplate(resourceName, `
		// 			resource "cloudavenue_vdc_group" "example_move_state" {
		// 				name = "move-state"
		// 				vdc_ids = [
		// 					cloudavenue_vdc.example_vdc_group_1.id,
		// 				]
		// 			}`),
		// 			Checks: []resource.TestCheckFunc{
		// 				resource.TestCheckResourceAttrWith("cloudavenue_vdc_group.example_move_state", "id", urn.TestIsType(urn.VDCGroup)),
		// 				resource.TestCheckResourceAttrSet("cloudavenue_vdc_group.example_move_state", "status"),
		// 				resource.TestCheckResourceAttrSet("cloudavenue_vdc_group.example_move_state", "type"),
		// 				resource.TestCheckResourceAttr("cloudavenue_vdc_group.example_move_state", "name", "move-state"),
		// 				resource.TestCheckResourceAttr("cloudavenue_vdc_group.example_move_state", "vdc_ids.#", "1"),
		// 			},
		// 		},
		// 		// ! Updates testing
		// 		Updates: []testsacc.TFConfig{
		// 			// Update description and add a new vdc_id
		// 			{
		// 				TFConfig: testsacc.GenerateFromTemplate(resourceName, `
		// 				resource "cloudavenue_vdcg" "example_move_state" {
		// 					name = "move-state"
		// 					vdc_ids = [
		// 						cloudavenue_vdc.example_vdc_group_1.id,
		// 					]
		// 				}

		// 				moved {
		// 					from = cloudavenue_vdc_group.example
		// 					to   = cloudavenue_vdcg.example
		// 				}`),
		// 				Checks: []resource.TestCheckFunc{
		// 					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.VDCGroup)),
		// 					resource.TestCheckResourceAttrSet(resourceName, "status"),
		// 					resource.TestCheckResourceAttrSet(resourceName, "type"),
		// 					resource.TestCheckResourceAttr(resourceName, "name", "move-state"),
		// 					resource.TestCheckResourceAttr(resourceName, "vdc_ids.#", "1"),
		// 				},
		// 			},
		// 		},
		// 	}
		// },
	}
}

func TestAccVDCGResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps:                    testsacc.GenerateTests(&VDCGResource{}),
	})
}
