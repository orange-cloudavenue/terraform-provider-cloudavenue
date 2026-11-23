// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package testsacc

import (
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// testIsURNType checks that the attribute value is a URN of the given type.
// An empty value is accepted, like cloudavenue-sdk-go/pkg/urn.TestIsType.
func testIsURNType(urnType urn.URN) resource.CheckResourceAttrWithFunc {
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
