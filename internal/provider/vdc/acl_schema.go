package vdc

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

/*
aclSchema
This function is used to create the superschema for the vDC ACL.
*/
func aclSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "Provides a Cloud Avenue vDC access control resource.",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "This can be used to share vDC across users and/or groups.",
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
			"everyone_access_level": acl.SuperSchema(true)["everyone_access_level"],
			"shared_with":           acl.SuperSchema(true)["shared_with"],
		},
	}
}
