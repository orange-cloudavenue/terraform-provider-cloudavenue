package diskparams

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func StorageProfileAttribute() schema.Attribute {
	return schema.StringAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Storage profile to override the VM default one.",
	}
}
