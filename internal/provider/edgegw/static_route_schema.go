package edgegw

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/orange-cloudavenue/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func staticRouteSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_static_route` resource allows you to create and manage static routes on an Edge Gateway.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_edgegateway_static_route` data source allows you to retrieve information about a static route on an Edge Gateway.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the Static Route.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
					Optional:            true,
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
						stringplanmodifier.UseStateForUnknown(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Static Route.",
					Required:            true,
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the Static Route.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"network_cidr": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The network CIDR of the Static Route. (e.g. 192.168.1.0/24)",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"next_hops": superschema.SuperSetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "A set of next hops to use within the static route",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
					Validators: []validator.Set{
						setvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"ip_address": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "IP address for next hop gateway IP Address for the Static Route.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								fstringvalidator.IsIP(),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"admin_distance": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Admin distance is used to choose which route to use when there are multiple routes for a specific network. The lower the admin distance, the higher the preference for the route.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
							Default:  int64default.StaticInt64(1),
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
