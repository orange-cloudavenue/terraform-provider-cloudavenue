/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package acl

import (
	"fmt"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
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
SuperSchema is the superschema for ACL.

If readOnly is false, returns ACL attributes with 3 rights : "ReadOnly", "Change", "FullControl".
If readOnly is true, returns ACL attributes with only "ReadOnly".
*/
func SuperSchema(readOnly bool) map[string]superschema.Attribute {
	validatorACL := stringvalidator.OneOf("ReadOnly", "Change", "FullControl")
	if readOnly {
		validatorACL = stringvalidator.OneOf("ReadOnly")
	}

	return map[string]superschema.Attribute{
		"everyone_access_level": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "Access level when the vApp is shared with everyone.",
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					validatorACL,
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
							validatorACL,
						},
					},
				},
			},
		},
	}
}

func SharedSetToAccessControl(_ *govcd.VCDClient, org *govcd.AdminOrg, input []SharedWithModel) ([]*govcdtypes.AccessSetting, []*SharedWithModel, error) {
	var output []*govcdtypes.AccessSetting
	var outputModel []*SharedWithModel

	for _, item := range input {
		var subjectHref string
		var subjectType string
		var subjectName string
		var oModel *SharedWithModel

		if !item.UserID.IsNull() && !item.UserID.IsUnknown() {
			userID := item.UserID.ValueString()
			user, err := org.GetUserById(userID, false)
			if err != nil {
				return nil, nil, fmt.Errorf("error retrieving user %s: %w", userID, err)
			}
			subjectHref = user.User.Href
			subjectType = user.User.Type
			subjectName = user.User.Name

			oModel = &SharedWithModel{
				UserID:      types.StringValue(urn.Normalize(urn.User, urn.ExtractUUID(subjectHref)).String()),
				SubjectName: types.StringValue(subjectName),
			}
		} else if !item.GroupID.IsNull() && !item.GroupID.IsUnknown() {
			groupID := item.GroupID.ValueString()
			group, err := org.GetGroupById(groupID, false)
			if err != nil {
				return nil, nil, fmt.Errorf("error retrieving group %s: %w", groupID, err)
			}
			subjectHref = group.Group.Href
			subjectType = group.Group.Type
			subjectName = group.Group.Name
			oModel = &SharedWithModel{
				GroupID:     types.StringValue(urn.Normalize(urn.Group, urn.ExtractUUID(subjectHref)).String()),
				SubjectName: types.StringValue(subjectName),
			}
		}

		accessLevel := item.AccessLevel.ValueString()

		output = append(output, &govcdtypes.AccessSetting{
			Subject: &govcdtypes.LocalSubject{
				HREF: subjectHref,
				Name: subjectName,
				Type: subjectType,
			},
			ExternalSubject: nil,
			AccessLevel:     accessLevel,
		})
		oModel.AccessLevel = types.StringValue(accessLevel)
		outputModel = append(outputModel, oModel)
	}
	return output, outputModel, nil
}

func AccessControlListToSharedSet(org *govcd.AdminOrg, input []*govcdtypes.AccessSetting) ([]SharedWithModel, error) {
	if org == nil {
		return nil, fmt.Errorf("admin organization is required to resolve ACL subjects")
	}

	var output []SharedWithModel

	for _, item := range input {
		if item == nil || item.Subject == nil {
			return nil, fmt.Errorf("access setting subject is missing")
		}
		if item.ExternalSubject != nil {
			return nil, fmt.Errorf("external subject '%s' is not supported", item.ExternalSubject.SubjectId)
		}

		o := SharedWithModel{}

		subjectUUID := urn.ExtractUUID(item.Subject.HREF)
		if subjectUUID == "" {
			return nil, fmt.Errorf("cannot extract subject UUID from href '%s' for item %s", item.Subject.HREF, item.Subject.Name)
		}

		if _, err := org.GetUserById(subjectUUID, false); err == nil {
			o.UserID = types.StringValue(urn.Normalize(urn.User, subjectUUID).String())
		} else if _, err := org.GetGroupById(subjectUUID, false); err == nil {
			o.GroupID = types.StringValue(urn.Normalize(urn.Group, subjectUUID).String())
		} else {
			switch item.Subject.Type {
			case govcdtypes.MimeAdminUser:
				o.UserID = types.StringValue(urn.Normalize(urn.User, subjectUUID).String())
			case govcdtypes.MimeAdminGroup:
				o.GroupID = types.StringValue(urn.Normalize(urn.Group, subjectUUID).String())
			default:
				return nil, fmt.Errorf("cannot resolve ACL subject '%s' (%s)", item.Subject.Name, item.Subject.HREF)
			}
		}
		o.AccessLevel = types.StringValue(item.AccessLevel)
		o.SubjectName = types.StringValue(item.Subject.Name)

		output = append(output, o)
	}
	return output, nil
}
