// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package vapp

import (
	"fmt"

	"github.com/orange-cloudavenue/common-go/urn"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdktypes "github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/acl"
)

// URN prefixes expected by the SDK v2 for local vApp ACL subjects.
// The provider keeps urn.User / urn.Group in the state.
const (
	sdkOrgUserURNPrefix  = "urn:vcloud:orguser:"
	sdkOrgGroupURNPrefix = "urn:vcloud:orggroup:"

	sdkACLKindUser  = "user"
	sdkACLKindGroup = "group"
)

type aclResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	VDC                 types.String `tfsdk:"vdc"`
	VAppID              types.String `tfsdk:"vapp_id"`
	VAppName            types.String `tfsdk:"vapp_name"`
	EveryoneAccessLevel types.String `tfsdk:"everyone_access_level"`
	SharedWith          types.Set    `tfsdk:"shared_with"`
}

// sharedWithToSDK converts the shared_with entries to SDK ACL subjects.
// Only local subjects are supported. User and group IDs are converted from the
// provider URN (urn:vcloud:user: / urn:vcloud:group:) to the SDK one (orguser / orggroup).
func sharedWithToSDK(input []acl.SharedWithModel) ([]sdktypes.ParamsSetVAppACLSubject, error) {
	output := make([]sdktypes.ParamsSetVAppACLSubject, 0, len(input))

	for _, item := range input {
		subject := sdktypes.ParamsSetVAppACLSubject{
			SubjectSource: sdktypes.VAppACLSubjectSourceLocal,
			AccessLevel:   item.AccessLevel.ValueString(),
		}

		switch {
		case !item.UserID.IsNull() && !item.UserID.IsUnknown():
			uuid := urn.ExtractUUID(item.UserID.ValueString())
			if uuid == "" {
				return nil, fmt.Errorf("invalid user_id %q: no UUID found", item.UserID.ValueString())
			}
			subject.SubjectKind = sdkACLKindUser
			subject.SubjectID = sdkOrgUserURNPrefix + uuid
		case !item.GroupID.IsNull() && !item.GroupID.IsUnknown():
			uuid := urn.ExtractUUID(item.GroupID.ValueString())
			if uuid == "" {
				return nil, fmt.Errorf("invalid group_id %q: no UUID found", item.GroupID.ValueString())
			}
			subject.SubjectKind = sdkACLKindGroup
			subject.SubjectID = sdkOrgGroupURNPrefix + uuid
		default:
			return nil, fmt.Errorf("shared_with entry needs a user_id or a group_id")
		}

		output = append(output, subject)
	}

	return output, nil
}

// sharedWithFromSDK converts SDK ACL subjects to shared_with entries.
// IDs are converted back to the provider URN to keep the state unchanged.
// External subjects are not supported, as before.
func sharedWithFromSDK(input []sdktypes.ModelVAppACLSubject) ([]acl.SharedWithModel, error) {
	output := make([]acl.SharedWithModel, 0, len(input))

	for _, subject := range input {
		if subject.SubjectSource != sdktypes.VAppACLSubjectSourceLocal {
			return nil, fmt.Errorf("external subject '%s' is not supported", subject.SubjectID)
		}

		uuid := urn.ExtractUUID(subject.SubjectID)
		if uuid == "" {
			return nil, fmt.Errorf("cannot extract subject UUID from id '%s' for item %s", subject.SubjectID, subject.SubjectName)
		}

		item := acl.SharedWithModel{
			AccessLevel: types.StringValue(subject.AccessLevel),
			SubjectName: types.StringValue(subject.SubjectName),
		}

		switch subject.SubjectKind {
		case sdkACLKindUser:
			item.UserID = types.StringValue(urn.Normalize(urn.User, uuid).String())
		case sdkACLKindGroup:
			item.GroupID = types.StringValue(urn.Normalize(urn.Group, uuid).String())
		default:
			return nil, fmt.Errorf("cannot resolve ACL subject '%s' (%s): unknown kind %q", subject.SubjectName, subject.SubjectID, subject.SubjectKind)
		}

		output = append(output, item)
	}

	return output, nil
}
