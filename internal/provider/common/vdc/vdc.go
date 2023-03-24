package vdc

import (
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	superschema "github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/schema"
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
*/
func SuperSchema() superschema.StringAttribute {
	return superschema.StringAttribute{
		Resource: &schemaR.StringAttribute{
			Optional: true,
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
				stringplanmodifier.RequiresReplace(),
			},
			MarkdownDescription: "(ForceNew) The name of vDC to use, optional if defined at provider level.",
		},
		DataSource: &schemaD.StringAttribute{
			Optional:            true,
			MarkdownDescription: "The name of vDC to use, optional if defined at provider level.",
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
	}

	return v, nil
}
