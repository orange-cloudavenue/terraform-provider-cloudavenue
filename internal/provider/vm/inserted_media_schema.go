package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

func vmInsertedMediaSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "The inserted_media resource resource for inserting or ejecting media (ISO) file for the VM. Create this resource for inserting the media, and destroy it for ejecting.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the inserted media. This is the vm Id where the media is inserted.",
			},
			"vdc":       vdc.Schema(),
			"vapp_id":   vapp.Schema()["vapp_id"],
			"vapp_name": vapp.Schema()["vapp_name"],
			"catalog": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the catalog where to find media file",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Media file name in catalog which will be inserted to VM",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "VM name where media will be inserted or ejected",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			// "eject_force": schema.BoolAttribute{ - Disable attributes - Issue referrer: vmware/go-vcloud-director#552
			//	Optional:            true,
			//	MarkdownDescription: "Allows to pass answer to question in vCD when ejecting from a VM which is powered on. True means 'Yes' as answer to question. Default is true",
			//	PlanModifiers: []planmodifier.Bool{
			//		boolpm.SetDefault(true),
			//	},
			// },
		},
	}
}
