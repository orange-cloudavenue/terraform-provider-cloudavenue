package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

const sizeInMBDescription = "The size of the disk in MB."

/*
SizeInMBAttribute

returns a schema.Attribute with a value.
*/
func SizeInMBAttribute() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: sizeInMBDescription,
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	}
}

// SizeInMBAttributeComputed returns a schema.Attribute with a computed value.
func SizeInMBAttributeComputed() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: sizeInMBDescription,
		Computed:            true,
	}
}

// SizeInMBAttributeRequired returns a schema.Attribute with a required value.
func SizeInMBAttributeRequired() schema.Attribute {
	return schema.Int64Attribute{
		MarkdownDescription: sizeInMBDescription,
		Required:            true,
	}
}
