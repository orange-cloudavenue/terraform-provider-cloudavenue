package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

const busNumberDescription = "The number of the `SCSI` or `IDE` controller itself."

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

// BusNumberAttributeComputed returns a schema.Attribute with a computed value.
func BusNumberAttributeComputed() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: busNumberDescription,
		Computed:            true,
	}
}

const unitNumberDescription = "The device number on the `SCSI` or `IDE` controller of the disk."

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
