package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

type aclResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
	VAppID              types.String `tfsdk:"vapp_id"`
	VAppName            types.String `tfsdk:"vapp_name"`
	EveryoneAccessLevel types.String `tfsdk:"everyone_access_level"`
	SharedWith          types.Set    `tfsdk:"shared_with"`
}

/*
aclSchema
This function is used to create the superschema for the vAPP ACL.
*/
func aclSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue Access Control structure for a vApp.",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "This can be used to create, update, and delete access control structures for a vApp.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the acl rule.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"vdc":                   vdc.SuperSchema(),
			"vapp_id":               vapp.SuperSchema()["vapp_id"],
			"vapp_name":             vapp.SuperSchema()["vapp_name"],
			"everyone_access_level": acl.SuperSchema(true)["everyone_access_level"],
			"shared_with":           acl.SuperSchema(true)["shared_with"],
		},
	}
}
