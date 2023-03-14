package vdc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type VDC struct {
	VDCOrVDCGroup client.VDCOrVDCGroupHandler
	*govcd.Org
}

/*
Schema

	Optional: true
	Computed: true
	RequiresReplace
	UseStateForUnknown
*/
func Schema() schema.StringAttribute {
	return schema.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplace(),
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "(ForceNew) The name of vDC to use, optional if defined at provider level.",
	}
}

/*
Init

If vDC is not defined at data source level, use the one defined at provider level.
*/
func Init(client *client.CloudAvenue, vdc types.String) (VDC, diag.Diagnostics) {
	d := diag.Diagnostics{}
	if vdc.IsNull() || vdc.IsUnknown() {
		if client.DefaultVDCExist() {
			vdc = types.StringValue(client.GetDefaultVDC())
		} else {
			d.AddError("Missing VDC", "VDC is required when not defined at provider level")
			return VDC{}, d
		}
	}
	// Request Org and vDC
	orgOut, vdcOut, err := client.GetOrgAndVDC(client.GetOrg(), vdc.ValueString())
	if err != nil {
		d.AddError("Error retrieving VDC", err.Error())
		return VDC{}, d
	}
	return VDC{VDCOrVDCGroup: vdcOut, Org: orgOut}, nil
}

// GetName give you the name of the vDC.
func (v VDC) GetName() string {
	name := ""
	switch vdc := v.VDCOrVDCGroup.(type) {
	case *govcd.Vdc:
		name = vdc.Vdc.Name
	case *govcd.VdcGroup:
		name = vdc.VdcGroup.Name
	}
	return name
}

// GetID give you the ID of the vDC.
func (v VDC) GetID() string {
	id := ""
	switch vdc := v.VDCOrVDCGroup.(type) {
	case *govcd.Vdc:
		id = vdc.Vdc.ID
	case *govcd.VdcGroup:
		id = vdc.VdcGroup.Id
	}
	return id
}

// GetOrg give you the Org of the vDC.
func (v VDC) GetOrg() *govcd.Org {
	return v.Org
}

// OrgID give you the ID of the Org of the vDC.
func (v VDC) GetOrgID() string {
	return v.Org.Org.ID
}

// GetVDC	return the vDC.
func (v VDC) GetVDC() (*govcd.Vdc, diag.Diagnostics) {
	d := diag.Diagnostics{}
	vdc, isVDC := v.VDCOrVDCGroup.(*govcd.Vdc)
	if !isVDC {
		d.AddError("error retrieving VDC", fmt.Sprintf("expected *govcd.Vdc type, have %T", v.VDCOrVDCGroup))
	}
	return vdc, d
}
