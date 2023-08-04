package catalog

import (
	"context"
	"fmt"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/adminorg"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type ACLModel struct {
	CatalogID           supertypes.StringValue    `tfsdk:"catalog_id"`
	CatalogName         supertypes.StringValue    `tfsdk:"catalog_name"`
	EveryoneAccessLevel supertypes.StringValue    `tfsdk:"everyone_access_level"`
	ID                  supertypes.StringValue    `tfsdk:"id"`
	SharedWithEveryone  supertypes.BoolValue      `tfsdk:"shared_with_everyone"`
	SharedWithUsers     supertypes.SetNestedValue `tfsdk:"shared_with_users"`
}

// * SharedWithUsers.
type ACLModelSharedWithUsers []ACLModelSharedWithUser

// * SharedWithUsers.
type ACLModelSharedWithUser struct {
	AccessLevel supertypes.StringValue `tfsdk:"access_level"`
	UserID      supertypes.StringValue `tfsdk:"user_id"`
}

func NewCatalogACL(t any) *ACLModel {
	switch x := t.(type) {
	case tfsdk.State: //nolint:dupl
		return &ACLModel{
			CatalogID:           supertypes.NewStringUnknown(),
			CatalogName:         supertypes.NewStringUnknown(),
			EveryoneAccessLevel: supertypes.NewStringNull(),
			ID:                  supertypes.NewStringUnknown(),
			SharedWithEveryone:  supertypes.NewBoolNull(),
			SharedWithUsers:     supertypes.NewSetNestedNull(x.Schema.GetAttributes()["shared_with_users"].GetType().(supertypes.SetNestedType).ElementType()),
		}
	case tfsdk.Plan: //nolint:dupl
		return &ACLModel{
			CatalogID:           supertypes.NewStringUnknown(),
			CatalogName:         supertypes.NewStringUnknown(),
			EveryoneAccessLevel: supertypes.NewStringNull(),
			ID:                  supertypes.NewStringUnknown(),
			SharedWithEveryone:  supertypes.NewBoolNull(),
			SharedWithUsers:     supertypes.NewSetNestedNull(x.Schema.GetAttributes()["shared_with_users"].GetType().(supertypes.SetNestedType).ElementType()),
		}
	case tfsdk.Config: //nolint:dupl
		return &ACLModel{
			CatalogID:           supertypes.NewStringUnknown(),
			CatalogName:         supertypes.NewStringUnknown(),
			EveryoneAccessLevel: supertypes.NewStringNull(),
			ID:                  supertypes.NewStringUnknown(),
			SharedWithEveryone:  supertypes.NewBoolNull(),
			SharedWithUsers:     supertypes.NewSetNestedNull(x.Schema.GetAttributes()["shared_with_users"].GetType().(supertypes.SetNestedType).ElementType()),
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *ACLModel) Copy() *ACLModel {
	x := &ACLModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetSharedWithUsers returns the value of the SharedWithUsers field.
func (rm *ACLModel) GetSharedWithUsers(ctx context.Context) (values ACLModelSharedWithUsers, diags diag.Diagnostics) {
	values = make(ACLModelSharedWithUsers, 0)
	d := rm.SharedWithUsers.Get(ctx, &values, false)
	return values, d
}

// * Custom Funcs

// ToControlAccessParams converts the ACLModel to govcdtypes.ControlAccessParams.
func (rm *ACLModel) ToControlAccessParams(ctx context.Context, adminOrg adminorg.AdminOrg) (values govcdtypes.ControlAccessParams, diags diag.Diagnostics) {
	// Init slice
	accessSettings := make([]*govcdtypes.AccessSetting, 0)

	if rm.SharedWithUsers.IsKnown() {
		sharedWithusers, d := rm.GetSharedWithUsers(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		for _, user := range sharedWithusers {
			userAPI, err := adminOrg.GetUserById(user.UserID.Get(), true)
			if err != nil {
				diags.AddError(fmt.Sprintf("error when retrieving user %s", user.UserID.Get()), err.Error())
				return
			}

			accessSettings = append(accessSettings, &govcdtypes.AccessSetting{
				Subject: &govcdtypes.LocalSubject{
					HREF: userAPI.User.Href,
					Type: userAPI.User.Type,
					Name: userAPI.User.Name,
				},
				AccessLevel:     user.AccessLevel.Get(),
				ExternalSubject: nil,
			})
		}
	}

	var accessSettingsList *govcdtypes.AccessSettingList
	if len(accessSettings) > 0 {
		accessSettingsList = &govcdtypes.AccessSettingList{
			AccessSetting: accessSettings,
		}
	} else {
		accessSettingsList = nil
	}

	return govcdtypes.ControlAccessParams{
		IsSharedToEveryone:  rm.SharedWithEveryone.Get(),
		EveryoneAccessLevel: rm.EveryoneAccessLevel.GetPtr(),
		AccessSettings:      accessSettingsList,
	}, diags
}
