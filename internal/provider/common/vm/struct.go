package vm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type VMResourceModel struct { //nolint:revive
	ID          types.String `tfsdk:"id"`
	VDC         types.String `tfsdk:"vdc"`
	Name        types.String `tfsdk:"name"`
	VappName    types.String `tfsdk:"vapp_name"`
	VappID      types.String `tfsdk:"vapp_id"`
	Description types.String `tfsdk:"description"`
	DeployOS    types.Object `tfsdk:"deploy_os"`
	State       types.Object `tfsdk:"state"`
	Resource    types.Object `tfsdk:"resource"`
	Settings    types.Object `tfsdk:"settings"`
}

type VMResourceModelAllStructs struct { //nolint:revive
	DeployOS *VMResourceModelDeployOS
	State    *VMResourceModelState
	Resource *VMResourceModelResource
	Settings *VMResourceModelSettings
}

// AllStructsFromPlan returns the values of all the attributes of the VMResourceModel, if set, as a *VMResourceModelAllStructs.
func (rm *VMResourceModel) AllStructsFromPlan(ctx context.Context) (allStructs *VMResourceModelAllStructs, diags diag.Diagnostics) {
	allStructs = &VMResourceModelAllStructs{}

	allStructs.DeployOS, diags = rm.DeployOSFromPlan(ctx)
	if diags.HasError() {
		return
	}

	allStructs.State, diags = rm.StateFromPlan(ctx)
	if diags.HasError() {
		return
	}

	allStructs.Resource, diags = rm.ResourceFromPlan(ctx)
	if diags.HasError() {
		return
	}

	allStructs.Settings, diags = rm.SettingsFromPlan(ctx)
	if diags.HasError() {
		return
	}

	return
}

// * DeployOS
// DeployOSFromPlan returns the value of the DeployOS attribute, if set, as a VMResourceModelDeployOS.
func (rm *VMResourceModel) DeployOSFromPlan(ctx context.Context) (deployOS *VMResourceModelDeployOS, diags diag.Diagnostics) {
	tflog.Info(ctx, "DeployOSFromPlan")

	if rm.DeployOS.IsNull() {
		return &VMResourceModelDeployOS{
			VappTemplateID:   types.StringNull(),
			VMNameInTemplate: types.StringNull(),
			BootImageID:      types.StringNull(),
			AcceptAllEulas:   types.BoolNull(),
		}, nil
	}

	deployOS = &VMResourceModelDeployOS{}

	diags.Append(rm.DeployOS.As(ctx, deployOS, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	return
}

// * State
// StateFromPlan returns the value of the State attribute, if set, as a VMResourceModelState.
func (rm *VMResourceModel) StateFromPlan(ctx context.Context) (state *VMResourceModelState, diags diag.Diagnostics) {
	tflog.Info(ctx, "StateFromPlan")

	if rm.State.IsNull() || rm.State.IsUnknown() {
		return &VMResourceModelState{
			PowerON: types.BoolNull(),
			Status:  types.StringNull(),
		}, nil
	}

	state = &VMResourceModelState{}

	diags.Append(rm.State.As(ctx, state, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	return
}

// * Resource
// ResourceFromPlan returns the value of the Resource attribute, if set, as a VMResourceModelResource.
func (rm *VMResourceModel) ResourceFromPlan(ctx context.Context) (resource *VMResourceModelResource, diags diag.Diagnostics) {
	tflog.Info(ctx, "ResourceFromPlan")

	networks := VMResourceModelResourceNetworks{}

	if rm.Resource.IsNull() || rm.Resource.IsUnknown() {
		return &VMResourceModelResource{
			CPUs:                types.Int64Null(),
			CPUsCores:           types.Int64Null(),
			CPUHotAddEnabled:    types.BoolNull(),
			Memory:              types.Int64Null(),
			MemoryHotAddEnabled: types.BoolNull(),
			Networks:            types.ListNull(networks.ObjectType()),
		}, nil
	}

	resource = &VMResourceModelResource{}

	diags.Append(rm.Resource.As(ctx, resource, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	return
}

// * Networks
// NetworksFromPlan returns the value of the Networks attribute, if set, as a VMResourceModelResourceNetworks.
func (r *VMResourceModelResource) NetworksFromPlan(ctx context.Context) (networks *VMResourceModelResourceNetworks, diags diag.Diagnostics) {
	tflog.Info(ctx, "NetworksFromPlan")

	networks = &VMResourceModelResourceNetworks{}

	if r.Networks.IsNull() || r.Networks.IsUnknown() {
		return
	}

	diags.Append(r.Networks.ElementsAs(ctx, networks, false)...)

	return
}

// * Settings
// SettingsFromPlan returns the value of the Settings attribute, if set, as a VMResourceModelSettings.
func (rm *VMResourceModel) SettingsFromPlan(ctx context.Context) (settings *VMResourceModelSettings, diags diag.Diagnostics) {
	tflog.Info(ctx, "SettingsFromPlan")

	gP := VMResourceModelSettingsGuestProperties{}
	custom := VMResourceModelSettingsCustomization{}

	if rm.Settings.IsNull() || rm.Settings.IsUnknown() {
		return &VMResourceModelSettings{
			Customization:                types.ObjectNull(custom.AttrTypes()),
			GuestProperties:              types.MapUnknown(gP.AttrType()),
			ExposeHardwareVirtualization: types.BoolNull(),
			OsType:                       types.StringNull(),
			StorageProfile:               types.StringNull(),
			AffinityRuleID:               types.StringNull(),
		}, nil
	}

	settings = &VMResourceModelSettings{}

	diags.Append(rm.Settings.As(ctx, settings, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)

	return
}

// * SettingsCustomization
// CustomizationFromPlan returns the value of the SettingsCustomization attribute, if set, as a VMResourceModelSettingsCustomization.
func (s *VMResourceModelSettings) CustomizationFromPlan(ctx context.Context) (customization *VMResourceModelSettingsCustomization, diags diag.Diagnostics) {
	tflog.Info(ctx, "CustomizationFromPlan")

	customization = &VMResourceModelSettingsCustomization{}

	if s.Customization.IsNull() || s.Customization.IsUnknown() {
		return customization, nil
	}

	diags.Append(s.Customization.As(ctx, customization, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: false,
	})...)

	return customization, diags
}

// * SettingsGuestProperties
// GuestPropertiesFromPlan returns the value of the SettingsGuestProperties attribute, if set, as a VMResourceModelSettingsGuestProperties.
func (s *VMResourceModelSettings) GuestPropertiesFromPlan(ctx context.Context, x types.Map) (guestProperties *VMResourceModelSettingsGuestProperties, diags diag.Diagnostics) {
	tflog.Info(ctx, "GuestPropertiesFromPlan")

	if s.GuestProperties.IsNull() || s.GuestProperties.IsUnknown() {
		return
	}

	guestProperties = &VMResourceModelSettingsGuestProperties{}

	diags = x.ElementsAs(ctx, guestProperties, false)

	return
}
