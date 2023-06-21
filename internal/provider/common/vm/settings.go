package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type VMResourceModelSettings struct { //nolint:revive
	ExposeHardwareVirtualization types.Bool   `tfsdk:"expose_hardware_virtualization"`
	OsType                       types.String `tfsdk:"os_type"`
	StorageProfile               types.String `tfsdk:"storage_profile"`
	GuestProperties              types.Map    `tfsdk:"guest_properties"`
	AffinityRuleID               types.String `tfsdk:"affinity_rule_id"`
	Customization                types.Object `tfsdk:"customization"`
}

// Equal returns true if the two VMResourceModelSettings are equal.
func (s *VMResourceModelSettings) Equal(other *VMResourceModelSettings) bool {
	return s.ExposeHardwareVirtualization.Equal(other.ExposeHardwareVirtualization) &&
		s.OsType.Equal(other.OsType) &&
		s.StorageProfile.Equal(other.StorageProfile) &&
		s.GuestProperties.Equal(other.GuestProperties) &&
		s.AffinityRuleID.Equal(other.AffinityRuleID) &&
		s.Customization.Equal(other.Customization)
}

// AttrTypes returns the types of the attributes of the Settings attribute.
func (s *VMResourceModelSettings) AttrTypes() map[string]attr.Type {
	return s.attrTypes(&VMResourceModelSettingsCustomization{}, &VMResourceModelSettingsGuestProperties{})
}

// attrTypes() returns the types of the attributes of the Settings attribute.
func (s *VMResourceModelSettings) attrTypes(customization *VMResourceModelSettingsCustomization, guestProperties *VMResourceModelSettingsGuestProperties) map[string]attr.Type {
	return map[string]attr.Type{
		"expose_hardware_virtualization": types.BoolType,
		"os_type":                        types.StringType,
		"storage_profile":                types.StringType,
		"guest_properties":               types.MapType{ElemType: guestProperties.AttrType()},
		"affinity_rule_id":               types.StringType,
		"customization":                  types.ObjectType{AttrTypes: customization.AttrTypes()},
	}
}

// toAttrValues() returns the values of the attributes of the Settings attribute.
func (s *VMResourceModelSettings) toAttrValues(_ context.Context) map[string]attr.Value {
	return map[string]attr.Value{
		"expose_hardware_virtualization": s.ExposeHardwareVirtualization,
		"os_type":                        s.OsType,
		"storage_profile":                s.StorageProfile,
		"guest_properties":               s.GuestProperties,
		"affinity_rule_id":               s.AffinityRuleID,
		"customization":                  s.Customization,
	}
}

// ToPlan returns the value of the Settings attribute, if set, as a types.Object.
func (s *VMResourceModelSettings) ToPlan(ctx context.Context) types.Object {
	if s == nil {
		return types.Object{}
	}

	var (
		customization   = &VMResourceModelSettingsCustomization{}
		guestProperties = &VMResourceModelSettingsGuestProperties{}
	)

	return types.ObjectValueMust(s.attrTypes(customization, guestProperties), s.toAttrValues(ctx))
}

// SettingsRead returns the value of the Settings attribute, if set, as a *VMResourceModelSettings.
func (v VM) SettingsRead(ctx context.Context, stateCustomization any) (settings *VMResourceModelSettings, err error) {
	guestProperties, err := v.GuestPropertiesRead()
	if err != nil {
		return nil, fmt.Errorf("unable to read guest properties: %w", err)
	}

	affinityRuleID, err := v.GetAffinityRuleIDOrDefault()
	if err != nil {
		return nil, fmt.Errorf("unable to read affinity rule ID: %w", err)
	}

	customization, err := v.CustomizationRead(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to read customization: %w", err)
	}

	switch custo := stateCustomization.(type) {
	case *VMResourceModelSettingsCustomization:
		customization.Force = custo.Force
	case attr.Value:
		x, ok := custo.(types.Object)
		if !ok {
			return nil, fmt.Errorf("unable to convert state customization to basetypes.ObjectType")
		}
		customization.Force = x.Attributes()["force"].(types.Bool)
	}

	return &VMResourceModelSettings{
		ExposeHardwareVirtualization: types.BoolValue(v.GetExposeHardwareVirtualization()),
		OsType:                       utils.StringValueOrNull(v.GetOSType()),
		StorageProfile:               utils.StringValueOrNull(v.GetStorageProfileName()),
		GuestProperties:              guestProperties.ToPlan(ctx),
		AffinityRuleID:               utils.StringValueOrNull(affinityRuleID),
		Customization:                customization.ToPlan(ctx),
	}, nil
}
