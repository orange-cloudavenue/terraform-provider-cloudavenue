package s3

import (
	"context"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

func credentialSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_credential` resource allows you to manage an access key and secret key for an S3 user.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the credential. ID is a username and 4 first characters of the access key. (e.g. `username-1234`).",
				},
			},
			"username": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The username is configured at the provider level.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
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
					MarkdownDescription: "If true, the API token will be saved in a file. Set this to true if you understand the security risks of using AccessKey/SecretKey files and agree to creating them. This setting is only used when creating a new AccessKey/SecretKey and available only one time.",
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
					MarkdownDescription: "If true, the API token will be printed in the console. Set this to true if you understand the security risks of using AccessKey/SecretKey and agree to creating them. This setting is only used when creating a new AccessKey/SecretKey and available only one time.",
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
					MarkdownDescription: "If true, the SecretKey will be saved in the terraform state. Set this to true if you understand the security risks of using AccessKey/SecretKey and agree to creating them. This setting is only used when creating a new API token and available only one time. \n\n !> **Warning:** This is a security risk and should only be used for testing purposes.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.RequiresReplace(),
						boolplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"access_key": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The Access Key.",
					Computed:            true,
					Sensitive:           true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"secret_key": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The Secret Key. Only Available if the `save_in_tfstate` is set to true.",
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
