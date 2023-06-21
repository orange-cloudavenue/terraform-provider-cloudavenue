package vm

import (
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
)

type securityTagResourceModel struct {
	Name  types.String `tfsdk:"id"`
	VMIDs types.Set    `tfsdk:"vm_ids"`
}

func securityTagSchema() superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The security_tag resource allows you to assign security tags to VMs.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Resource: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "ID is the name of the security tag.",
					Validators: []validator.String{
						stringvalidator.LengthBetween(1, 129),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"vm_ids": superschema.SetAttribute{
				Resource: &schemaR.SetAttribute{
					Required:            true,
					MarkdownDescription: "The IDs of the VMs to attach to the security tag.",
					ElementType:         types.StringType,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
						setvalidator.ValueStringsAre(fstringvalidator.IsURN()),
					},
				},
			},
		},
	}
}
