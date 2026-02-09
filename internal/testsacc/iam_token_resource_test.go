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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

const testAccTokenResourceConfig = `
resource "cloudavenue_iam_token" "example" {
	name            = "example"

	save_in_tfstate = true
	save_in_file    = true
	print_token     = true
}
`

const testAccTokenResourceConfigUpdate = `
resource "cloudavenue_iam_token" "example" {
	name            = "exampleUpdated"

	save_in_tfstate = true
	save_in_file    = true
	print_token     = true
}
`

func TestAccTokenResource(t *testing.T) {
	resourceName := "cloudavenue_iam_token.example"

	t.Cleanup(deleteFile("token.json", t))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// Apply test
				Config: testAccTokenResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "id", urn.TestIsType(urn.Token)),
					resource.TestCheckResourceAttr(resourceName, "name", "example"),
					resource.TestCheckResourceAttr(resourceName, "save_in_tfstate", "true"),
					resource.TestCheckResourceAttr(resourceName, "save_in_file", "true"),
					resource.TestCheckResourceAttr(resourceName, "print_token", "true"),
					resource.TestCheckResourceAttr(resourceName, "file_name", "token.json"),
					resource.TestCheckResourceAttrSet(resourceName, "token"),
					testCheckFileExists("token.json"),
				),
			},
			{
				// Update test
				// Any change generates replacement
				Config:             testAccTokenResourceConfigUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
