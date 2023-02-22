package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func SizeInMBAttribute() schema.Attribute {
	return schema.Int64Attribute{
		Required:            true,
		MarkdownDescription: "The size of the disk in MB.",
	}
}
