package s3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
)

func bucketCorsConfigurationSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_cors_configuration` resource allows you to manage the [Cross-Origin Resource Sharing](https://docs.aws.amazon.com/AmazonS3/latest/userguide/cors.html) configuration of an S3 bucket.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_cors_configuration` data source allows you to retrieve information about an S3 bucket's [Cross-Origin Resource Sharing](https://docs.aws.amazon.com/AmazonS3/latest/userguide/cors.html) configuration.",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Update: true,
					Delete: true,
					Read:   true,
				},
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the bucket",
				},
			},
			"bucket": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The name of the bucket.",
					Required:            true,
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"cors_rules": superschema.SuperSetNestedAttribute{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "Set of origins and methods (cross-origin access that you want to allow).",
					Validators: []validator.Set{
						setvalidator.SizeBetween(1, 100),
					},
				},
				Resource: &schemaR.SetNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"allowed_headers": superschema.SuperSetAttribute{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "Set of Headers that are specified in the Access-Control-Request-Headers header.",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"allowed_methods": superschema.SuperSetAttribute{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "Set of HTTP methods that you allow the origin to execute.",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Required: true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(stringvalidator.OneOf("GET", "PUT", "HEAD", "POST", "DELETE")),
							},
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"allowed_origins": superschema.SuperSetAttribute{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "Set of origins you want customers to be able to access the bucket from.",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Required: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"expose_headers": superschema.SuperSetAttribute{
						Common: &schemaR.SetAttribute{
							MarkdownDescription: "Set of headers in the response that you want customers to be able to access from their applications (for example, from a JavaScript XMLHttpRequest object).",
							ElementType:         supertypes.StringType{},
						},
						Resource: &schemaR.SetAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetAttribute{
							Computed: true,
						},
					},
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Unique identifier for the rule.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"max_age_seconds": superschema.SuperInt64Attribute{
						Common: &schemaR.Int64Attribute{
							MarkdownDescription: "Time in seconds that your browser is to cache the preflight response for the specified resource.",
						},
						Resource: &schemaR.Int64Attribute{
							Optional: true,
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
