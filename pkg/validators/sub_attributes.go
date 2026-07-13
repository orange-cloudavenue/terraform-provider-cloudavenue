/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package validators provides custom validation logic for Terraform resources.
package validators

import (
	"context"

	"github.com/orange-cloudavenue/terraform-plugin-framework-validators/listvalidator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SubAttributesNullIfAppIDValuesIsNotOne returns a validator that allows sub_attributes only if app_id.values contains exactly one value.
func SubAttributesNullIfAppIDValuesIsNotOne() validator.List {
	desc := listvalidator.NullIfAttributeMatchesDescription{
		Description:         "`sub_attributes` is only allowed if `app_id.values` contains exactly one value.",
		MarkdownDescription: "`app_id.sub_attributes` is only allowed if `app_id.values` contains exactly one value.",
	}
	return listvalidator.NullIfAttributeMatchesWithDescription(
		path.MatchRoot("app_id").AtName("values"),
		func(ctx context.Context, value attr.Value) (bool, diag.Diagnostics) {
			if value.IsNull() || value.IsUnknown() {
				return false, nil
			}

			list, ok := value.(types.List)
			if !ok {
				return false, nil
			}

			return list.Elements() != nil && len(list.Elements()) == 1, nil
		},
		desc,
	)
}
