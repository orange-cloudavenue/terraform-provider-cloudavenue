package vdc

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type VDC struct {
	*client.VDC
	org *client.Org
}

/*
Schema

	Optional: true
	Computed: true
	RequiresReplace
	UseStateForUnknown
*/
func Schema() schemaR.StringAttribute {
	return schemaR.StringAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
			stringplanmodifier.RequiresReplace(),
		},
		MarkdownDescription: "(ForceNew) The name of vDC to use, optional if defined at provider level.",
	}
}

/*
SuperSchema

	For the resource :
	Optional: true
	Computed: true
	RequiresReplace
	UseStateForUnknown

	For the data source :
	Optional: true
	Computed: true
*/
func SuperSchema() superschema.StringAttribute {
	return superschema.StringAttribute{
		Common: &schemaR.StringAttribute{
			MarkdownDescription: "The name of vDC to use, optional if defined at provider level.",
			Optional:            true,
			Computed:            true,
		},
		Resource: &schemaR.StringAttribute{
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
		},
	}
}

/*
Init

If vDC is not defined at data source level, use the one defined at provider level.
*/
func Init(c *client.CloudAvenue, vdc types.String) (VDC, diag.Diagnostics) {
	var (
		d    = diag.Diagnostics{}
		opts = make([]client.GetVDCOpts, 0)
		v    = VDC{}

		err error
	)

	v.org, err = c.GetOrg()
	if err != nil {
		d.AddError("Unable to get ORG", err.Error())
		return VDC{}, d
	}

	if !vdc.IsNull() && !vdc.IsUnknown() {
		opts = append(opts, client.WithVDCName(vdc.ValueString()))
	}

	v.VDC, err = c.GetVDC(opts...)
	if err != nil {
		d.AddError("Unable to get VDC", err.Error())
		return VDC{}, d
	}

	return v, nil
}
