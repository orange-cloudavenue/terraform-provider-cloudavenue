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

package s3

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"
)

type UserDataSourceModel struct {
	ID          supertypes.StringValue `tfsdk:"id"`
	Username    supertypes.StringValue `tfsdk:"user_name"`
	UserID      supertypes.StringValue `tfsdk:"user_id"`
	FullName    supertypes.StringValue `tfsdk:"full_name"`
	CanonicalID supertypes.StringValue `tfsdk:"canonical_id"`
}
