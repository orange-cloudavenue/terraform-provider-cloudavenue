package iam

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func tokenSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_iam_token` resource allows you to manage API tokens for a specific user. \n\n !> **Warning:** This resource print or store the token in clear text.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The unique ID of the API token for a specific user.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"name": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The unique name of the API token for a specific user.",
					Required:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"file_name": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the file to store the API token.",
					Optional:            true,
					Computed:            true,
					Default:             stringdefault.StaticString("token.json"),
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"save_in_file": superschema.SuperBoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "If true, the API token will be saved in a file. Set this to true if you understand the security risks of using API token files and agree to creating them. This setting is only used when creating a new API token and available only one time.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
					},
				},
			},
			"print_token": superschema.SuperBoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "If true, the API token will be printed in the console. Set this to true if you understand the security risks of using API token and agree to creating them. This setting is only used when creating a new API token and available only one time.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
					},
				},
			},
			"save_in_tfstate": superschema.SuperBoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "If true, the API token will be saved in the terraform state. Set this to true if you understand the security risks of using API token and agree to creating them. This setting is only used when creating a new API token and available only one time. \n\n !> **Warning:** This is a security risk and should only be used for testing purposes.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
						boolplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"token": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The API token for a specific user. Only Available if the `save_in_tfstate` is set to true.",
					Computed:            true,
					Sensitive:           true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
	}
}
