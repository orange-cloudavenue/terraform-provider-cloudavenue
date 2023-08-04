package catalog

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	fsetvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/setvalidator"
	fstringvalidator "github.com/FrangipaneTeam/terraform-plugin-framework-validators/stringvalidator"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/uuid"
)

func aclSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_catalog_acl` resource allows you to manage catalog ACLs in Cloud Avenue.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_catalog_acl` data source allows you to retrieve information about an ACL in Cloud Avenue.",
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The ID is same as the ID of the catalog.",
					Computed:            true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"catalog_id": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the catalog.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name"), path.MatchRoot("catalog_id")),
						fstringvalidator.IsURN(),
						fstringvalidator.PrefixContains(uuid.Catalog.String()),
					},
				},
			},
			"catalog_name": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The Name of the catalog.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.ExactlyOneOf(path.MatchRoot("catalog_name"), path.MatchRoot("catalog_id")),
					},
				},
			},
			"shared_with_everyone": superschema.SuperBoolAttribute{
				Resource: &schemaR.BoolAttribute{
					MarkdownDescription: "Whether the Catalog is shared with everyone in your organization with right `ReadOnly`.",
					Optional:            true,
					Computed:            true,
					Default:             booldefault.StaticBool(false),
				},
			},
			"everyone_access_level": superschema.SuperStringAttribute{
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "Access level when the Catalog is shared with everyone",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("ReadOnly", "Change", "FullControl"),
						fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRoot("shared_with_everyone"), []attr.Value{types.BoolValue(true)}),
						fstringvalidator.NullIfAttributeIsOneOf(path.MatchRoot("shared_with_everyone"), []attr.Value{types.BoolValue(false)}),
					},
				},
			},
			"shared_with_users": superschema.SuperSetNestedAttribute{
				Resource: &schemaR.SetNestedAttribute{
					MarkdownDescription: "The list of users with whom the Catalog is shared.",
					Optional:            true,
					Validators: []validator.Set{
						fsetvalidator.NullIfAttributeIsOneOf(path.MatchRoot("shared_with_everyone"), []attr.Value{types.BoolValue(true)}),
					},
				},
				Attributes: superschema.Attributes{
					"user_id": superschema.SuperStringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The ID of the user to which we are sharing.",
							Required:            true,
							Validators: []validator.String{
								fstringvalidator.IsURN(),
								fstringvalidator.PrefixContains(uuid.User.String()),
							},
						},
					},
					"access_level": superschema.SuperStringAttribute{
						Resource: &schemaR.StringAttribute{
							MarkdownDescription: "The access level for the user to which we are sharing.",
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("ReadOnly"),
							Validators: []validator.String{
								stringvalidator.OneOf("ReadOnly", "Change", "FullControl"),
							},
						},
					},
				},
			},
		},
	}
}
