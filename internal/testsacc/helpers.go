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
	"fmt"

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// ToValidate returns a CheckResourceAttrWithFunc that validates a string using the given validator name.
// If the value is empty, an error is returned indicating the value is not valid for the validator.
// Otherwise, the value is validated using the provided validator name via validators.New().Var().
func ToValidate(validatorToCheck string) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if value == "" {
			return fmt.Errorf("empty value, is not a valid %s", validatorToCheck)
		}

		return validators.New().Var(value, validatorToCheck)
	}
}
