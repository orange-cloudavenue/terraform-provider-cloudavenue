package diskparams

import (
	"strings"

	fstringplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	busTypeIDE  = busType{name: "ide", code: "5", subtype: "ide"}                      // Bus type IDE
	busTypeSATA = busType{name: "sata", code: "20", subtype: "vmware.sata.ahci"}       // Bus type SATA
	busTypeSCSI = busType{name: "scsi", code: "6", subtype: "lsilogicsas"}             // Bus type SCSI
	busTypeNVME = busType{name: "nvme", code: "20", subtype: "vmware.nvme.controller"} // Bus type NVME
)

type busType struct {
	name    string
	code    string
	subtype string
}

func (b busType) Name() string {
	return strings.ToUpper(b.name)
}

func (b busType) SubType() string {
	return b.subtype
}

func (b busType) Code() string {
	return b.code
}

func GetBusTypeByCode(code, subtype string) busType {
	switch code {
	case busTypeSATA.code:
		// SATA and NVME have the same code
		// Using the subtype to differentiate them
		switch subtype {
		case busTypeNVME.subtype:
			return busTypeNVME
		default:
			return busTypeSATA
		}
	case busTypeSCSI.code:
		return busTypeSCSI
	default:
		return busTypeSATA
	}
}

func GetBusTypeByName(name string) busType {
	switch strings.ToLower(name) {
	case busTypeSATA.name:
		return busTypeSATA
	case busTypeSCSI.name:
		return busTypeSCSI
	case busTypeNVME.name:
		return busTypeNVME
	default:
		return busTypeSATA
	}
}

var listOfBusTypes = []string{busTypeIDE.Name(), busTypeSATA.Name(), busTypeSCSI.Name(), busTypeNVME.Name()}

const busTypeDescription = "The type of disk controller. Possible values: `scsi`, `sata` or `nvme`. Default value is `scsi`."

/*
BusTypeAttribute

returns a schema.Attribute with a value.

This is Optional and has a default value of busTypeSCSI.String().
*/
func BusTypeAttribute() schema.Attribute {
	return schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: busTypeDescription,
		Validators: []validator.String{
			stringvalidator.OneOf(listOfBusTypes...),
		},
		PlanModifiers: []planmodifier.String{
			fstringplanmodifier.SetDefault(busTypeSCSI.Name()),
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

// BusTypeAttributeComputed returns a schema.Attribute with a computed value.
func BusTypeAttributeComputed() schema.Attribute {
	return schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: busTypeDescription,
	}
}

// BusTypeAttributeRequired returns a schema.Attribute with a required value.
func BusTypeAttributeRequired() schema.Attribute {
	return schema.StringAttribute{
		Required:            true,
		MarkdownDescription: busTypeDescription,
		Validators: []validator.String{
			stringvalidator.OneOf(listOfBusTypes...),
		},
	}
}
