package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type busType string

func (b busType) String() string {
	return string(b)
}

const (
	busTypeIDE         busType = "ide"         // Bus type IDE
	busTypeParallel    busType = "parallel"    // Bus type Parallel (LSI Logic Parallel SCSI)
	busTypeSAS         busType = "sas"         // Bus type SAS (LSI Logic SAS SCSI)
	busTypeParavirtual busType = "paravirtual" // Bus type Paravirtual (Paravirtual SCSI)
	busTypeSATA        busType = "sata"        // Bus type SATA
	busTypeNVME        busType = "nvme"        // Bus type NVME
)

var busTypes = []string{busTypeIDE.String(), busTypeParallel.String(), busTypeSAS.String(), busTypeParavirtual.String(), busTypeSATA.String(), busTypeNVME.String()}

func BusTypeAttribute() schema.Attribute {
	return schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The type of disk controller. Possible values: `ide`, `parallel` (LSI Logic Parallel SCSI), `sas` (LSI Logic SAS SCSI), `paravirtual` (Paravirtual SCSI), `sata`, `nvme`.",
		Validators: []validator.String{
			stringvalidator.OneOf(busTypes...),
		},
	}
}
