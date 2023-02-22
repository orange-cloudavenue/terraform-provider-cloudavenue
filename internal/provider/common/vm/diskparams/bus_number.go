package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func BusNumberAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Required: true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.RequiresReplace(),
		},
		MarkdownDescription: "The number of the `SCSI` or `IDE` controller itself.",
	}
}
