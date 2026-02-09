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

package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

const sizeInMBDescription = "The size of the disk in MB."

/*
SizeInMBAttribute

returns a schema.Attribute with a value.
*/
func SizeInMBAttribute() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: sizeInMBDescription,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	}
}

// SizeInMBAttributeComputed returns a schema.Attribute with a computed value.
func SizeInMBAttributeComputed() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: sizeInMBDescription,
		Computed:            true,
	}
}

// SizeInMBAttributeRequired returns a schema.Attribute with a required value.
func SizeInMBAttributeRequired() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: sizeInMBDescription,
		Required:            true,
	}
}
