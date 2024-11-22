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
	// * FrangipaneTeam Custom Validators.
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"
)

// How to use types generator:
// 1. Define the schema in the file internal/provider/alb/virtual_service_schema.go
// 2. Add the resource or data source to the file internal/provider/provider_resources.go or internal/provider/provider_data_sources.go respectively
// 3. Launch the following command to generate golang structs for the schema:
// go run ./cmd/types-generator/*.go -file internal/provider/alb/virtual_service_schema.go -resource cloudavenue_alb_virtual_service -is-resource.
func virtualServiceSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "Provides a resource to manage ALB Virtual services in CloudAvenue. A virtual service advertises an IP address and ports to the external world and listens for client traffic. When a virtual service receives traffic, it directs it to members in ALB Pool.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_alb_virtual_service` data source allows you to retrieve information about an ...",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the virtual service.",
				},
			},
			"name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the ALB Virtual Service.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"edge_gateway_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the edge gateway on which the ALB Virtual Service is to be created.",
					Optional:            true,
					Computed:            true,
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
					Optional:            true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"enabled": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Defines if the ALB Virtual Service is enabled.",
					Optional:            true,
					Default:             booldefault.StaticBool(true),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
			"pool_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the ALB Server Pool associated.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("pool_name"), path.MatchRoot("pool_id")),
					},
				},
			},
			"pool_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the ALB Server Pool associated.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("pool_name"), path.MatchRoot("pool_id")),
					},
				},
			},
			"service_engine_group_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the service Engine Group.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"virtual_ip": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The virtual IP address of the ALB Virtual Service.",
					Optional:            true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_type": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The type of the service port.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.OneOf("HTTP", "HTTPS", "L4_TCP", "L4_UDP", "L4_TLS"),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"certificate_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the certificate.",
					Optional:            true,
				},
				Resource: &schemaR.StringAttribute{
					Validators: []validator.String{
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("service_type"), []attr.Value{types.StringValue("L4_TLS"), types.StringValue("HTTPS")}),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"service_port": superschema.SuperSetNestedAttributeOf[VirtualServiceModelServicePort]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The service port of the ALB Virtual Service.",
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"port_start": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The start port of the service port range or exact port number if `port_end`is not set.",
							Required:            true,
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"port_end": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "The end port of the service port range. If not specified, only the `port_start` value is used.",
							Optional:            true,
						},
						DataSource: &schemaD.Int64Attribute{
							Computed: true,
						},
					},
					"port_type": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type of the service port. The value is `UDP_FAST_PATH` must be used if you choose `L4_UDP. in the `service_type` attribute." + `
							A TCP/UDP fast path profile does not proxy TCP connections. It directly connects clients to the destination server and translates the destination virtual service address of the client with the IP address of the chosen destination server.
							`,
						},
						Resource: &schemaR.StringAttribute{
							Default:  stringdefault.StaticString("TCP_PROXY"),
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("TCP_PROXY", "TCP_FAST_PATH", "UDP_FAST_PATH"),
								fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("service_type"), []attr.Value{types.StringValue("L4_UDP")}),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"port_ssl": superschema.SuperBoolAttribute{
						Common: &schemaR.BoolAttribute{
							MarkdownDescription: "Defines if the service port is SSL enabled.",
							Optional:            true,
							Default:             booldefault.StaticBool(false),
						},
						DataSource: &schemaD.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
			"preserve_client_ip": superschema.SuperBoolAttribute{
				Common: &schemaR.BoolAttribute{
					MarkdownDescription: "Defines if the client IP address is preserved (proxy mode transparent).",
					Optional:            true,
					Default:             booldefault.StaticBool(false),
				},
				DataSource: &schemaD.BoolAttribute{
					Computed: true,
				},
			},
		},
	}
}
