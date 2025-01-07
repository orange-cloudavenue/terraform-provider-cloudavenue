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

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type BucketModel struct {
	ID         supertypes.StringValue `tfsdk:"id"`
	Name       supertypes.StringValue `tfsdk:"name"`
	ObjectLock supertypes.BoolValue   `tfsdk:"object_lock"`
	Endpoint   supertypes.StringValue `tfsdk:"endpoint"`
}

func (rm *BucketModel) Copy() *BucketModel {
	x := &BucketModel{}
	utils.ModelCopy(rm, x)
	return x
}
