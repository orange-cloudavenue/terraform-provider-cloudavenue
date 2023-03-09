package acl

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SharedWithModel struct {
	UserID      types.String `tfsdk:"user_id"`
	GroupID     types.String `tfsdk:"group_id"`
	AccessLevel types.String `tfsdk:"access_level"`
	SubjectName types.String `tfsdk:"subject_name"`
}

var SharedWithModelAttrTypes = map[string]attr.Type{
	"user_id":      types.StringType,
	"group_id":     types.StringType,
	"access_level": types.StringType,
	"subject_name": types.StringType,
}

/*
Schema

If readOnly is false, returns ACL attributes with 3 rights : "ReadOnly", "Change", "FullControl".
If readOnly is true, returns ACL attributes with only "ReadOnly".
*/
func Schema(readOnly bool) map[string]schema.Attribute {
	validatorACL := stringvalidator.OneOf("ReadOnly", "Change", "FullControl")
	descriptionACL := "one of `ReadOnly`, `Change`, `FullControl`"
	if readOnly {
		validatorACL = stringvalidator.OneOf("ReadOnly")
		descriptionACL = "only `ReadOnly` is available"
	}

	return map[string]schema.Attribute{
		"everyone_access_level": schema.StringAttribute{
			MarkdownDescription: "Access level when the vApp is shared with everyone (" + descriptionACL + ").\n" +
				"Only one of `everyone_access_level` or `shared_with` can be set.",
			Optional: true,
			Validators: []validator.String{
				validatorACL,
				stringvalidator.ExactlyOneOf(path.MatchRoot("shared_with"), path.MatchRoot("everyone_access_level")),
			},
		},
		"shared_with": schema.SetNestedAttribute{
			MarkdownDescription: "One or more blocks defining a subject to which we are sharing.\n" +
				"Only one of `everyone_access_level` or `shared_with` can be set.",
			Optional: true,
			Validators: []validator.Set{
				setvalidator.ExactlyOneOf(path.MatchRoot("everyone_access_level"), path.MatchRoot("shared_with")),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"user_id": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "ID of the user to which we are sharing.\n" +
							"Only one of `group_id` or `user_id` can be set.",
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("group_id")),
						},
					},
					"group_id": schema.StringAttribute{
						Optional: true,
						MarkdownDescription: "ID of the group to which we are sharing.\n" +
							"Only one of `group_id` or `user_id` can be set.",
						Validators: []validator.String{
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("user_id")),
						},
					},
					"subject_name": schema.StringAttribute{
						MarkdownDescription: "Name of the subject (group or user) with which we are sharing",
						Computed:            true,
					},
					"access_level": schema.StringAttribute{
						MarkdownDescription: "Access level for the user or group to which we are sharing (" + descriptionACL + ").",
						Required:            true,
						Validators: []validator.String{
							validatorACL,
						},
					},
				},
			},
		},
	}
}
