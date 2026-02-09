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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
)

const busNumberDescription = "The number of the controller itself."

/*
BusNumberAttribute

returns a schema.Attribute with a value.
if value is not set, the api try to compute value automaticaly.
*/
func BusNumberAttribute() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: busNumberDescription,
		Optional:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AtLeast(0),
			int64validator.AtMost(3),
		},
	}
}

// BusNumberAttributeRequired returns a schema.Attribute with a required value.
func BusNumberAttributeRequired() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: busNumberDescription,
		Required:            true,
		Validators: []validator.Int64{
			int64validator.AtLeast(0),
			int64validator.AtMost(3),
		},
	}
}

// BusNumberAttributeComputed returns a schema.Attribute with a computed value.
func BusNumberAttributeComputed() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: busNumberDescription,
		Computed:            true,
	}
}

const unitNumberDescription = "The device number on the controller of the disk."

/*
UnitNumberAttribute

returns a schema.Attribute with a value.
if value is not set, the api try to compute value automaticaly.
*/
func UnitNumberAttribute() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: unitNumberDescription,
		Optional:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
		Validators: []validator.Int64{
			int64validator.AtLeast(0),
			int64validator.AtMost(15),
		},
	}
}

// UnitNumberAttributeRequired returns a schema.Attribute with a required value.
func UnitNumberAttributeRequired() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: unitNumberDescription,
		Required:            true,
		Validators: []validator.Int64{
			int64validator.AtLeast(0),
			int64validator.AtMost(15),
		},
	}
}

// UnitNumberAttributeComputed returns a schema.Attribute with a computed value.
func UnitNumberAttributeComputed() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: unitNumberDescription,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
	}
}

// Max BusNumber is 4 (0,1,2,3)
// Max UnitNumber is 16 (0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15)
/*
Compute BusNumber and UnitNumber

if busNumber is not set, the api try to compute value automaticaly.
*/
func ComputeBusAndUnitNumber(disks []*govcdtypes.DiskSettings) (busNumber, unitNumber int) {
	busNumber = 0
	unitNumber = 0
	for _, disk := range disks {
		if disk.BusNumber > busNumber {
			busNumber = disk.BusNumber
		}
		if disk.UnitNumber > unitNumber {
			unitNumber = disk.UnitNumber
		}
	}
	if unitNumber == 15 {
		busNumber++
		unitNumber = 0
	} else {
		unitNumber++
	}
	return busNumber, unitNumber
}
