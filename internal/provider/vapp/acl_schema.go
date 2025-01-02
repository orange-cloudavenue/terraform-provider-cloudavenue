package vapp

import (
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

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
			"vdc":       vdc.SuperSchema(),
			"vapp_id":   vapp.SuperSchema()["vapp_id"],
			"vapp_name": vapp.SuperSchema()["vapp_name"],
			"everyone_access_level": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Access level when the vApp is shared with everyone.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.OneOf("ReadOnly", "Change", "FullControl"),
						stringvalidator.ExactlyOneOf(path.MatchRoot("shared_with"), path.MatchRoot("everyone_access_level")),
					},
				},
			},
			"shared_with": superschema.SetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "One or more blocks defining the subjects with whom we are sharing.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.ExactlyOneOf(path.MatchRoot("everyone_access_level"), path.MatchRoot("shared_with")),
					},
				},
				Attributes: map[string]superschema.Attribute{
					"user_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "ID of the user with whom we are sharing.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("group_id")),
							},
						},
					},
					"group_id": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "ID of the group with whom we are sharing.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("user_id")),
							},
						},
					},
					"subject_name": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Name of the subject (group or user) with whom we are sharing",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
					},
					"access_level": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Access level for the user or group with whom we are sharing.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("ReadOnly", "Change", "FullControl"),
							},
						},
					},
				},
			},
		},
	}
}
