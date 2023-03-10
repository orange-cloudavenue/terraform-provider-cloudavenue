// Package vapp provides common functionality for vApp resources.
package vapp

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

const (
	ErrVAppNotFound = "vApp not found"
)

type VApp struct {
	*govcd.VApp
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
func Schema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"vapp_id": schema.StringAttribute{
			MarkdownDescription: "(ForceNew) ID of the vApp. Required if `vapp_name` is not set.",
			Optional:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Validators: []validator.String{
				stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
			},
		},
		"vapp_name": schema.StringAttribute{
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
Init

Get vApp name or vApp ID.
*/
func Init(client *client.CloudAvenue, vdc vdc.VDC, vappID, vappName types.String) (vapp VApp, d diag.Diagnostics) {
	var vappNameID string
	if vappID.IsNull() || vappID.IsUnknown() {
		vappNameID = vappName.ValueString()
	} else {
		vappNameID = vappID.ValueString()
	}

	// Request vApp
	vappOut, err := vdc.GetVAppByNameOrId(vappNameID, true)
	if err != nil {
		if errors.Is(err, govcd.ErrorEntityNotFound) {
			d.AddError(ErrVAppNotFound, err.Error())
			return
		}
		d.AddError("Error retrieving vApp", err.Error())
		return
	}
	return VApp{VApp: vappOut, vdc: vdc}, nil
}

// LockParentVApp locks the parent vApp.
func (v VApp) LockParentVApp(ctx context.Context) (d diag.Diagnostics) {
	if v.vdc.GetName() == "" || v.GetName() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "vDC: "+v.vdc.GetName()+" vApp: "+v.GetName())
		return
	}
	key := fmt.Sprintf("vdc:%s|vapp:%s", v.vdc.GetName(), v.GetName())
	vcdMutexKV.KvLock(ctx, key)
	return
}

// UnlockParentVApp unlocks the parent vApp.
func (v VApp) UnlockParentVApp(ctx context.Context) (d diag.Diagnostics) {
	if v.vdc.GetName() == "" || v.GetName() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "vDC: "+v.vdc.GetName()+" vApp: "+v.GetName())
		return
	}
	key := fmt.Sprintf("vdc:%s|vapp:%s", v.vdc.GetName(), v.GetName())
	vcdMutexKV.KvUnlock(ctx, key)
	return
}

// GetName give you the name of the vApp.
func (v VApp) GetName() string {
	return v.VApp.VApp.Name
}

// GetID give you the ID of the vApp.
func (v VApp) GetID() string {
	return v.VApp.VApp.ID
}

// GetStatusCode give you the status code of the vApp.
func (v VApp) GetStatusCode() int {
	return v.VApp.VApp.Status
}

// GetHREF give you the HREF of the vApp.
func (v VApp) GetHREF() string {
	return v.VApp.VApp.HREF
}

// GetDescription give you the status code of the vApp.
func (v VApp) GetDescription() string {
	return v.VApp.VApp.Description
}
