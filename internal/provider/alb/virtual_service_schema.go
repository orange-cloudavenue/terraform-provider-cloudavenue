package alb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

func virtualServiceSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a resource to manage ALB Virtual services for particular NSX-T Edge Gateway. A virtual service advertises an IP address and ports to the external world and listens for client traffic. When a virtual service receives traffic, it directs it to members in ALB Pool.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a data source to read ALB Virtual services for particular NSX-T Edge Gateway. A virtual service advertises an IP address and ports to the external world and listens for client traffic. When a virtual service receives traffic, it directs it to members in ALB Pool.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the load balancer virtual service.",
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the ALB Virtual Service.",
					Required:            true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the edge gateway on which the ALB Virtual Service is to be created.",
					Optional:            true,
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"edge_gateway_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the edge gateway on which the ALB Virtual Service is to be created.",
					Optional:            true,
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplaceIfConfigured(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
			},
			"description": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The description of the ALB Virtual Service.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Defines if the ALB Virtual Service is enabled.",
				},
				Resource: &schemaR.BoolAttribute{
					Optional: true,
					Computed: true,
					Default:  booldefault.StaticBool(true),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"pool_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the ALB Server Pool associated.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("pool_name"), path.MatchRoot("pool_id")),
					},
				},
			},
			"pool_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the ALB Server Pool associated.",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("pool_name"), path.MatchRoot("pool_id")),
					},
				},
			},
			"service_engine_group_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the service Engine Group (Take the first one if not specified).",
					Computed:            true,
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
				},
			},
			"virtual_ip": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The virtual IP address of the ALB Virtual Service.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the service. The different modes that the NSX Advanced Load Balancer supports for handling TCP traffic and various parameters that can be tuned for optimization of the TCP traffic are also detailed here.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						fstringvalidator.OneOfWithDescription(func() []fstringvalidator.OneOfWithDescriptionValues {
							var values []fstringvalidator.OneOfWithDescriptionValues
							for _, v := range v1.EdgeGatewayALBVirtualServiceModelApplicationProfiles {
								values = append(values, fstringvalidator.OneOfWithDescriptionValues{Value: string(v.Value), Description: v.Description})
							}
							return values
						}()...),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"certificate_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the certificate. The certificate must be uploaded to the NSX Advanced Load Balancer before it can be used. The certificate MUST'NT be expired.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Validators: []validator.String{
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("service_type"), []attr.Value{types.StringValue("L4_TLS"), types.StringValue("HTTPS")}),
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("service_ports").AtAnyListIndex().AtName("port_ssl"), []attr.Value{types.BoolValue(true)}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_ports": superschema.SuperListNestedAttributeOf[VirtualServiceModelServicePort]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "The service port of the ALB Virtual Service. The service port is the port on which the virtual service listens for client traffic.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"port_start": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The start port of the service port range or exact port number if `port_end`is not set.",
						},
						Resource: &schemaR.Int64Attribute{
							Required: true,
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"port_end": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The end port of the service port range. If not specified, only the `port_start` value is used.",
							Computed:            true,
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
						},
					},
					"port_type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the service port. The different modes that the NSX Advanced Load Balancer supports for handling TCP traffic and various parameters that can be tuned for optimization of the TCP traffic are also detailed here.",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Default:  stringdefault.StaticString("TCP_PROXY"),
							Optional: true,
							Validators: []validator.String{
								fstringvalidator.OneOfWithDescription(
									fstringvalidator.OneOfWithDescriptionValues{Value: "TCP_PROXY", Description: "The TCP proxy terminates client connections to the virtual service, processes the payload, and then opens a new TCP connection to the destination server. Any application data from the client that is destined for a server is forwarded to that server over the new server-side TCP connection. Separating (or proxying) the client-to-server connections enables the NSX Advanced Load Balancer to provide enhanced security, such as TCP protocol sanitization and denial of service (DoS) mitigation."},
									fstringvalidator.OneOfWithDescriptionValues{Value: "TCP_FAST_PATH", Description: "A TCP fast path profile does not proxy TCP connections. It directly connects clients to the destination server and translates the destination virtual service address of the client with the IP address of the chosen destination server. The source IP address of the client can be NATed to the IP address of the SE."},
									fstringvalidator.OneOfWithDescriptionValues{Value: "UDP_FAST_PATH", Description: "Advanced Load Balancer translates the client’s destination virtual service address to the destination server and writes the source IP address of the client to the address of the SE, when forwarding the packet to the server. This ensures that server response traffic traverses symmetrically through the original SE."},
								),
								fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("service_type"), []attr.Value{types.StringValue("L4")}),
								fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("service_type"), []attr.Value{types.StringValue("HTTP"), types.StringValue("HTTPS"), types.StringValue("L4_TLS")}),
							},
						},
					},
					"port_ssl": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Defines if the service port is SSL enabled.",
							Computed:            true,
						},
						Resource: &schemaR.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
			},
			// Not supported in cloudavenue (need edge gateway with mode transparent enabled)
			// "preserve_client_ip": superschema.SuperBoolAttribute{
			// 	Common: &schemaR.BoolAttribute{
			// 		MarkdownDescription: "Defines if the client IP address is preserved (proxy mode transparent).",
			// 		Optional:            true,
			// 		Default:             booldefault.StaticBool(false),
			// 	},
			// 	DataSource: &schemaD.BoolAttribute{
			// 		Computed: true,
			// 	},
			// },
		},
	}
}
