package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func bucketWebsiteConfigurationSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_website_configuration` resource allows you to manage a configuration to host static websites content. [For more information](https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteHosting.html)",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_website_configuration` data source allows you to retrieve information about a configuration of static websites content. [For more information](https://docs.aws.amazon.com/AmazonS3/latest/userguide/WebsiteHosting.html)",
		},
		Attributes: superschema.Attributes{
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
					MarkdownDescription: "The ID of the bucket website. This is the same as the bucket name",
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
			"index_document": superschema.SuperSingleNestedAttributeOf[BucketWebsiteConfigurationModelIndexDocument]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The name of the index document.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"suffix": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The suffix of the index document.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRoot("redirect_all_requests_to")),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"error_document": superschema.SuperSingleNestedAttributeOf[BucketWebsiteConfigurationModelErrorDocument]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "The name of the error document for the website.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"key": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The key of the error document.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.ConflictsWith(path.MatchRoot("redirect_all_requests_to")),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"redirect_all_requests_to": superschema.SuperSingleNestedAttributeOf[BucketWebsiteConfigurationModelRedirectAllRequestsTo]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "Redirect behavior for every request to this bucket's website endpoint. [For more information](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration#redirect_all_requests_to)",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRoot("index_document"), path.MatchRoot("error_document"), path.MatchRoot("routing_rules")),
					},
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"hostname": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Name of the host where requests will be redirected.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"protocol": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Protocol to use when redirecting requests. The default is the protocol that is used in the original request.",
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf(s3.Protocol_Values()...),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"routing_rules": superschema.SuperSetNestedAttributeOf[BucketWebsiteConfigurationModelRoutingRule]{
				Common: &schemaR.SetNestedAttribute{
					MarkdownDescription: "Rules that define when a redirect is applied and the redirect behavior. [For more information](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration#routing_rule)",
				},
				Resource: &schemaR.SetNestedAttribute{
					Optional: true,
					Validators: []validator.Set{
						setvalidator.ConflictsWith(path.MatchRoot("redirect_all_requests_to")),
					},
				},
				DataSource: &schemaD.SetNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"condition": superschema.SuperSingleNestedAttributeOf[BucketWebsiteConfigurationModelCondition]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Configuration block for describing a condition that must be met for the specified redirect to apply. [For more information](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration#condition)",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"key_prefix_equals": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The object key name prefix when the redirect is applied. For example, to redirect requests for `ExamplePage.html`, the key prefix will be `ExamplePage.html`. To redirect request for all pages with the prefix `docs/`, the key prefix will be `docs`, which identifies all objects in the docs/ folder.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(path.MatchRoot("redirect_all_requests_to")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"http_error_code_returned_equals": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The HTTP error code when the redirect is applied. In the event of an error, if the error code equals this value, then the specified redirect is applied. Required when parent element `key_prefix_equals` is specified and parent element `http_redirect` is not.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(path.MatchRoot("redirect_all_requests_to")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
						},
					},
					"redirect": superschema.SuperSingleNestedAttributeOf[BucketWebsiteConfigurationModelRedirect]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Configuration block for redirecting all requests to another host instead of the original host. [For more information](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket_website_configuration#redirect)",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"hostname": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Name of the host where requests will be redirected.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"http_redirect_code": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "HTTP redirect code to use on the response.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("301", "302"),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"protocol": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Protocol to use when redirecting requests. The default is the protocol that is used in the original request.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.OneOf("http", "https"),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"replace_key_prefix_with": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Object key prefix to use in the redirect request. For example, to redirect requests for all pages with prefix `docs/` (objects in the `docs/` folder) to `documents/`, you can set a condition block with `key_prefix_equals` set to `docs/` and in the redirect set `replace_key_prefix_with` to `/documents`.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("replace_key_with")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"replace_key_with": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Specific object key to use in the redirect request. For example, redirect request to `error.html`.",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("replace_key_prefix_with")),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
			"website_endpoint": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The website endpoint.",
					Computed:            true,
				},
			},
		},
	}
}
