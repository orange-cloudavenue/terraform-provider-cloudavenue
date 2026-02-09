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

package vm

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type DisksModel struct {
	Disks    supertypes.ListNestedValue `tfsdk:"disks"`
	VDC      supertypes.StringValue     `tfsdk:"vdc"`
	ID       supertypes.StringValue     `tfsdk:"id"`
	VAppID   supertypes.StringValue     `tfsdk:"vapp_id"`
	VAppName supertypes.StringValue     `tfsdk:"vapp_name"`
	VMID     supertypes.StringValue     `tfsdk:"vm_id"`
	VMName   supertypes.StringValue     `tfsdk:"vm_name"`
}

// * Disks.
type DisksModelDisks []DisksModelDisk

// * Disks.
type DisksModelDisk struct {
	ID             supertypes.StringValue `tfsdk:"id"`
	IsDetachable   supertypes.BoolValue   `tfsdk:"is_detachable"`
	Name           supertypes.StringValue `tfsdk:"name"`
	SizeInMb       supertypes.Int64Value  `tfsdk:"size_in_mb"`
	StorageProfile supertypes.StringValue `tfsdk:"storage_profile"`
}

func (rm *DisksModel) Copy() *DisksModel {
	x := &DisksModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetDisks returns the value of the Disks field.
func (rm *DisksModel) GetDisks(ctx context.Context) (values DisksModelDisks, diags diag.Diagnostics) {
	values = make(DisksModelDisks, 0)
	d := rm.Disks.Get(ctx, &values, false)
	return values, d
}
