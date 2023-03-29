package publicip

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

type publicIPResourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	ID              types.String   `tfsdk:"id"`
	PublicIP        types.String   `tfsdk:"public_ip"`
	EdgeGatewayName types.String   `tfsdk:"edge_gateway_name"`
	EdgeGatewayID   types.String   `tfsdk:"edge_gateway_id"`
}

func publicIPSchema() superschema.Schema {
	return superschema.Schema{
		Common: superschema.SchemaDetails{
			MarkdownDescription: "This allows you to",
		},
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "manage a Public IP in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "retrieve information about a Public IP in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Public IP.",
					Computed:            true,
				},
			},
			"public_ip": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Public IP Address.",
					Computed:            true,
				},
			},
			"edge_gateway_name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the Edge Gateway.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"edge_gateway_id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the Edge Gateway.",
				},
				Resource: &schemaR.StringAttribute{
					Optional: true,
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("edge_gateway_name"), path.MatchRoot("edge_gateway_id")),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
