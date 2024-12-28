package diskparams

import (
	"strings"

	fstringplanmodifier "github.com/orange-cloudavenue/terraform-plugin-framework-planmodifiers/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var (
	BusTypeIDE  = BusType{name: "ide", code: "5", subtype: "ide"}                      // Bus type IDE
	BusTypeSATA = BusType{name: "sata", code: "20", subtype: "vmware.sata.ahci"}       // Bus type SATA
	BusTypeSCSI = BusType{name: "scsi", code: "6", subtype: "lsilogicsas"}             // Bus type SCSI
	BusTypeNVME = BusType{name: "nvme", code: "20", subtype: "vmware.nvme.controller"} // Bus type NVME
)

type BusType struct {
	name    string
	code    string
	subtype string
}

func (b BusType) Name() string {
	return strings.ToUpper(b.name)
}

func (b BusType) SubType() string {
	return b.subtype
}

func (b BusType) Code() string {
	return b.code
}

func GetBusTypeByCode(code, subtype string) BusType {
	switch code {
	case BusTypeSATA.code:
		// SATA and NVME have the same code
		// Using the subtype to differentiate them
		switch subtype {
		case BusTypeNVME.subtype:
			return BusTypeNVME
		default:
			return BusTypeSATA
		}
	case BusTypeSCSI.code:
		return BusTypeSCSI
	default:
		return BusTypeSATA
	}
}

func GetBusTypeByName(name string) BusType {
	switch strings.ToLower(name) {
	case BusTypeSATA.name:
		return BusTypeSATA
	case BusTypeSCSI.name:
		return BusTypeSCSI
	case BusTypeNVME.name:
		return BusTypeNVME
	default:
		return BusTypeSATA
	}
}

var ListOfBusTypes = []string{BusTypeIDE.Name(), BusTypeSATA.Name(), BusTypeSCSI.Name(), BusTypeNVME.Name()}

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
			stringvalidator.OneOf(ListOfBusTypes...),
		},
		PlanModifiers: []planmodifier.String{
			fstringplanmodifier.SetDefault(BusTypeSCSI.Name()),
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
			stringvalidator.OneOf(ListOfBusTypes...),
		},
	}
}
