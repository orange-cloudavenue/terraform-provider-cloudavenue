/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package helpers

import (
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Special functions for the terraform provider test.
// TestIsType returns true if the URN is of the specified type.
func TestIsType(urnType urn.URN) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if value == "" {
			return nil
		}

		if !urn.URN(value).IsType(urnType) {
			return fmt.Errorf("urn %s is not of type %s", value, urnType)
		}
		return nil
	}
}
