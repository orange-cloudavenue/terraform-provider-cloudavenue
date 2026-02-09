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

// Package vapp provides common functionality for vApp resources.
package vapp

import (
	"context"
	"errors"
	"fmt"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

const (
	ErrVAppNotFound = "vApp not found"
)

const (
	SchemaVappID   = "vapp_id"
	SchemaVappName = "vapp_name"
)

type VAPP struct {
	*v1.VAPP
	vdc vdc.VDC
}

var (
	// ErrVAppRefEmpty is returned when a vApp reference is missing information.
	ErrVAppRefEmpty  = errors.New("missing information in vapp ref")
	vcdMutexKV       = mutex.NewKV()
	DiagVAppNotFound = diag.NewErrorDiagnostic(ErrVAppNotFound, govcd.ErrorEntityNotFound.Error())
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
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_name"), path.MatchRoot("vapp_id")),
				},
				Optional: true,
			},
			Resource: &schemaR.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
		"vapp_name": superschema.StringAttribute{
			Common: &schemaR.StringAttribute{
				MarkdownDescription: "Name of the vApp.",
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.MatchRoot("vapp_id"), path.MatchRoot("vapp_name")),
				},
				Optional: true,
			},
			Resource: &schemaR.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			DataSource: &schemaD.StringAttribute{
				Computed: true,
			},
		},
	}
}

/*
Init

Get vApp name or vApp ID.
If mustExist is false, returns a warning for non existent vApp.
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
			d.Append(diag.Diagnostics{DiagVAppNotFound}...)
			return vapp, d
		}
		d.AddError("Error retrieving vApp", err.Error())
		return vapp, d
	}
	return VAPP{VAPP: vappOut, vdc: vdc}, nil
}

/*
Create

Create vApp and return VAPP struct.
*/
func Create(vdc vdc.VDC, vappName, description string) (vapp VAPP, d diag.Diagnostics) {
	vappOut, err := vdc.CreateVAPP(vappName, description)
	if err != nil {
		d.AddError("Error creating vApp", err.Error())
		return vapp, d
	}
	return VAPP{VAPP: vappOut, vdc: vdc}, nil
}

// LockVAPP locks the parent vApp.
func (v VAPP) LockVAPP(ctx context.Context) (d diag.Diagnostics) {
	if v.vdc.GetName() == "" || v.GetName() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "vDC: "+v.vdc.GetName()+" vApp: "+v.GetName())
		return d
	}
	key := fmt.Sprintf("vdc:%s|vapp:%s", v.vdc.GetName(), v.GetName())
	vcdMutexKV.KvLock(ctx, key)
	return d
}

// UnlockVAPP unlocks the parent vApp.
func (v VAPP) UnlockVAPP(ctx context.Context) (d diag.Diagnostics) {
	if v.vdc.GetName() == "" || v.GetName() == "" || ctx == nil {
		d.AddError("Incorrect lock args", "vDC: "+v.vdc.GetName()+" vApp: "+v.GetName())
		return d
	}
	key := fmt.Sprintf("vdc:%s|vapp:%s", v.vdc.GetName(), v.GetName())
	vcdMutexKV.KvUnlock(ctx, key)
	return d
}
