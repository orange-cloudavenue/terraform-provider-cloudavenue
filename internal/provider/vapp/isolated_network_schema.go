package vapp

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	fboolplanmodifier "github.com/FrangipaneTeam/terraform-plugin-framework-planmodifiers/boolplanmodifier"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vapp"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/vdc"
)

func isolatedNetworkSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Provides capability to attach an existing Org VDC Network to a vApp and toggle network features.",
		Attributes: map[string]schema.Attribute{
			"vdc":       vdc.Schema(),
			"vapp_id":   vapp.Schema()["vapp_id"],
			"vapp_name": vapp.Schema()["vapp_name"],
			"guest_vlan_allowed": schema.BoolAttribute{
				MarkdownDescription: "True if Network allows guest VLAN. Default to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					fboolplanmodifier.SetDefault(false),
				},
			},
			"retain_ip_mac_enabled": schema.BoolAttribute{
				MarkdownDescription: "Specifies whether the network resources such as IP/MAC of router will be retained across deployments. Default to `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					fboolplanmodifier.SetDefault(false),
				},
			},
		},
	}
}
