// Package vapp provides common functionality for vApp resources.
package vapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vm"
)

const (
	ErrVAppNotFound = "vApp not found"
)

const (
	SchemaVappID   = "vapp_id"
	SchemaVappName = "vapp_name"
)

type VAPP struct {
	*client.VAPP
	vdc vdc.VDC
}

var (
	// ErrVAppRefEmpty is returned when a vApp reference is missing information.
	ErrVAppRefEmpty = errors.New("missing information in vapp ref")
	vcdMutexKV      = mutex.NewKV()
)

/*
Schema

Return the schema for vapp_id and vapp_name with MarkdownDescription, Validators and PlanModifiers.
*/
func Schema() map[string]schemaR.Attribute {
	return map[string]schemaR.Attribute{
		"vapp_id": schemaR.StringAttribute{
			MarkdownDescription: "(ForceNew) ID of the vApp. Required if `vapp_name` is not set.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
			},
		},
		"vapp_name": schemaR.StringAttribute{
			MarkdownDescription: "(ForceNew) Name of the vApp. Required if `vapp_id` is not set.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_id"), path.MatchRoot("vapp_name")),
			},
		},
	}
}

/*
SuperSchema

Return the superschema for vapp_id and vapp_name with MarkdownDescription, Validators and PlanModifiers.
*/
func SuperSchema() map[string]superschema.Attribute {
	return map[string]superschema.Attribute{
		"vapp_id": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "ID of the vApp.",
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
				},
			},
		},
		"vapp_name": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "Name of the vApp.",
			},
			Resource: &schemaR.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_id"), path.MatchRoot("vapp_name")),
				},
			},
		},
	}
}

/*
Init

Get vApp name or vApp ID.
*/
func Init(_ *client.CloudAvenue, vdc vdc.VDC, vappID, vappName types.String) (vapp VAPP, d diag.Diagnostics) {
	vappNameID := vappID.ValueString()
	if vappID.IsNull() || vappID.IsUnknown() {
		vappNameID = vappName.ValueString()
	}

	// Request vApp
	vappOut, err := vdc.GetVAPP(vappNameID, true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			d.AddError(ErrVAppNotFound, err.Error())
			return
		}
		d.AddError("Error retrieving vApp", err.Error())
		return
	}
	return VAPP{VAPP: vappOut, vdc: vdc}, nil
}

type GetVMOpts struct {
	ID   types.String
	Name types.String
}

// vmIDOrName returns the ID or name of the VM.
func (v GetVMOpts) vmIDOrName() string {
	if v.ID.IsNull() || v.ID.IsUnknown() {
		return v.Name.ValueString()
	}
	return v.ID.ValueString()
}

// GetVM returns a VM from a vApp.
func (v VAPP) GetVM(vmInfo GetVMOpts, refresh bool) (vm.VM, diag.Diagnostics) {
	var d diag.Diagnostics

	vmOut, err := v.GetVMByNameOrId(vmInfo.vmIDOrName(), refresh)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			d.AddError("VM not found", err.Error())
			return vm.VM{}, nil
		}
		d.AddError("Error retrieving VM", err.Error())
		return vm.VM{}, nil
	}
	return vm.VM{VM: &client.VM{VM: vmOut}}, nil
}

// LockVAPP locks the parent vApp.
func (v VAPP) LockVAPP(ctx context.Context) (d diag.Diagnostics) {
	if v.vdc.GetName() == "" || v.GetName() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "vDC: "+v.vdc.GetName()+" vApp: "+v.GetName())
		return
	}
	key := fmt.Sprintf("vdc:%s|vapp:%s", v.vdc.GetName(), v.GetName())
	vcdMutexKV.KvLock(ctx, key)
	return
}

// UnlockVAPP unlocks the parent vApp.
func (v VAPP) UnlockVAPP(ctx context.Context) (d diag.Diagnostics) {
	if v.vdc.GetName() == "" || v.GetName() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "vDC: "+v.vdc.GetName()+" vApp: "+v.GetName())
		return
	}
	key := fmt.Sprintf("vdc:%s|vapp:%s", v.vdc.GetName(), v.GetName())
	vcdMutexKV.KvUnlock(ctx, key)
	return
}

// CreateVMWithTemplate.
func (v VAPP) CreateVMWithTemplate() (vm vm.VM, d diag.Diagnostics) {
	// vmFromTemplateParams := &govcdtypes.ReComposeVAppParams{
	// 	Ovf:              govcdtypes.XMLNamespaceOVF,
	// 	Xsi:              govcdtypes.XMLNamespaceXSI,
	// 	Xmlns:            govcdtypes.XMLNamespaceVCloud,
	// 	AllEULAsAccepted: v.Plan.AcceptAllEulas.ValueBool(),
	// 	Name:             vapp.VApp.Name,
	// 	PowerOn:          false, // VM will be powered on after all configuration is done
	// 	SourcedItem: &govcdtypes.SourcedCompositionItemParam{
	// 		Source: &govcdtypes.Reference{
	// 			HREF: vmTemplate.VAppTemplate.HREF,
	// 			Name: v.Plan.VMName.ValueString(), // This VM name defines the VM name after creation
	// 		},
	// 		VMGeneralParams: &govcdtypes.VMGeneralParams{
	// 			Description: v.Plan.Description.ValueString(),
	// 		},
	// 		InstantiationParams: &govcdtypes.InstantiationParams{
	// 			// If a MAC address is specified for NIC - it does not get set with this call,
	// 			// therefore an additional `vm.UpdateNetworkConnectionSection` is required.
	// 			NetworkConnectionSection: &networkConnectionSection,
	// 		},
	// 		ComputePolicy:  vmComputePolicy,
	// 		StorageProfile: storageProfilePtr,
	// 	},
	// }

	return
}

// CreateVMWithBootImage
