/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vm

import (
	"context"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VMResourceModelSettingsCustomization struct { //nolint:revive
	Force                          types.Bool   `tfsdk:"force"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	ChangeSID                      types.Bool   `tfsdk:"change_sid"`
	AllowLocalAdminPassword        types.Bool   `tfsdk:"allow_local_admin_password"`
	MustChangePasswordOnFirstLogin types.Bool   `tfsdk:"must_change_password_on_first_login"`
	AdminPassword                  types.String `tfsdk:"admin_password"`
	AutoGeneratePassword           types.Bool   `tfsdk:"auto_generate_password"`
	NumberOfAutoLogons             types.Int64  `tfsdk:"number_of_auto_logons"`
	JoinDomain                     types.Bool   `tfsdk:"join_domain"`
	JoinOrgDomain                  types.Bool   `tfsdk:"join_org_domain"`
	JoinDomainName                 types.String `tfsdk:"join_domain_name"`
	JoinDomainUser                 types.String `tfsdk:"join_domain_user"`
	JoinDomainPassword             types.String `tfsdk:"join_domain_password"`
	JoinDomainAccountOU            types.String `tfsdk:"join_domain_account_ou"`
	InitScript                     types.String `tfsdk:"init_script"`
	Hostname                       types.String `tfsdk:"hostname"`
}

// attrTypes() returns the types of the attributes of the SettingsCustomization attribute.
func (s *VMResourceModelSettingsCustomization) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"force":                               types.BoolType,
		"enabled":                             types.BoolType,
		"change_sid":                          types.BoolType,
		"allow_local_admin_password":          types.BoolType,
		"must_change_password_on_first_login": types.BoolType,
		"admin_password":                      types.StringType,
		"auto_generate_password":              types.BoolType,
		"number_of_auto_logons":               types.Int64Type,
		"join_domain":                         types.BoolType,
		"join_org_domain":                     types.BoolType,
		"join_domain_name":                    types.StringType,
		"join_domain_user":                    types.StringType,
		"join_domain_password":                types.StringType,
		"join_domain_account_ou":              types.StringType,
		"init_script":                         types.StringType,
		"hostname":                            types.StringType,
	}
}

// toAttrValues() returns the values of the attributes of the SettingsCustomization attribute.
func (s *VMResourceModelSettingsCustomization) toAttrValues() map[string]attr.Value {
	return map[string]attr.Value{
		"force":                               s.Force,
		"enabled":                             s.Enabled,
		"change_sid":                          s.ChangeSID,
		"allow_local_admin_password":          s.AllowLocalAdminPassword,
		"must_change_password_on_first_login": s.MustChangePasswordOnFirstLogin,
		"admin_password":                      s.AdminPassword,
		"auto_generate_password":              s.AutoGeneratePassword,
		"number_of_auto_logons":               s.NumberOfAutoLogons,
		"join_domain":                         s.JoinDomain,
		"join_org_domain":                     s.JoinOrgDomain,
		"join_domain_name":                    s.JoinDomainName,
		"join_domain_user":                    s.JoinDomainUser,
		"join_domain_password":                s.JoinDomainPassword,
		"join_domain_account_ou":              s.JoinDomainAccountOU,
		"init_script":                         s.InitScript,
		"hostname":                            s.Hostname,
	}
}

// ToPlan returns the value of the SettingsCustomization attribute, if set, as a types.Object.
func (s *VMResourceModelSettingsCustomization) ToPlan(_ context.Context) types.Object {
	if s == nil {
		return types.ObjectNull(s.AttrTypes())
	}

	return types.ObjectValueMust(s.AttrTypes(), s.toAttrValues())
}

// CustomizationRead reads the customization fields from a VM.
func (v *VM) CustomizationRead(_ context.Context) (*VMResourceModelSettingsCustomization, error) {
	customization, err := v.GetCustomization()
	if err != nil {
		return nil, err
	}

	return &VMResourceModelSettingsCustomization{
		Force:                          types.BoolValue(false),
		Enabled:                        utils.BoolPtrValueOrNull(customization.Enabled),
		ChangeSID:                      utils.BoolPtrValueOrNull(customization.ChangeSid),
		AllowLocalAdminPassword:        utils.BoolPtrValueOrNull(customization.AdminPasswordEnabled),
		MustChangePasswordOnFirstLogin: utils.BoolPtrValueOrNull(customization.ResetPasswordRequired),
		AdminPassword:                  utils.StringValueOrNull(customization.AdminPassword),
		AutoGeneratePassword:           utils.BoolPtrValueOrNull(customization.AdminPasswordAuto),
		NumberOfAutoLogons:             types.Int64Value(int64(customization.AdminAutoLogonCount)),
		JoinDomain:                     utils.BoolPtrValueOrNull(customization.JoinDomainEnabled),
		JoinOrgDomain:                  utils.BoolPtrValueOrNull(customization.UseOrgSettings),
		JoinDomainName:                 utils.StringValueOrNull(customization.DomainName),
		JoinDomainUser:                 utils.StringValueOrNull(customization.DomainUserName),
		JoinDomainPassword:             utils.StringValueOrNull(customization.DomainUserPassword),
		JoinDomainAccountOU:            utils.StringValueOrNull(customization.MachineObjectOU),
		InitScript:                     utils.StringValueOrNull(customization.CustomizationScript),
		Hostname:                       utils.StringValueOrNull(customization.ComputerName),
	}, nil
}

// GetCustomizationSection returns the value of the SettingsCustomization attribute, if set, as a *types.CustomizationSection.
func (s *VMResourceModelSettingsCustomization) GetCustomizationSection(vmName string) *govcdtypes.GuestCustomizationSection {
	if s.Hostname.ValueString() != "" {
		vmName = s.Hostname.ValueString()
	}

	return &govcdtypes.GuestCustomizationSection{
		ComputerName:        vmName,
		Enabled:             s.Enabled.ValueBoolPointer(),
		ChangeSid:           s.ChangeSID.ValueBoolPointer(),
		CustomizationScript: s.InitScript.ValueString(),

		AdminPasswordEnabled: s.AllowLocalAdminPassword.ValueBoolPointer(),
		AdminPasswordAuto:    s.AutoGeneratePassword.ValueBoolPointer(),
		AdminPassword:        s.AdminPassword.ValueString(),

		ResetPasswordRequired: s.MustChangePasswordOnFirstLogin.ValueBoolPointer(),

		AdminAutoLogonEnabled: utils.TakeBoolPointer(s.NumberOfAutoLogons.ValueInt64() > 0),
		AdminAutoLogonCount:   int(s.NumberOfAutoLogons.ValueInt64()),

		JoinDomainEnabled:  s.JoinDomain.ValueBoolPointer(),
		UseOrgSettings:     s.JoinOrgDomain.ValueBoolPointer(),
		DomainName:         s.JoinDomainName.ValueString(),
		DomainUserName:     s.JoinDomainUser.ValueString(),
		DomainUserPassword: s.JoinDomainPassword.ValueString(),
		MachineObjectOU:    s.JoinDomainAccountOU.ValueString(),
	}
}
