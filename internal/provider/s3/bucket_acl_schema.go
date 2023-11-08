package s3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

func bucketACLSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_acl` resource allows you to manage S3 (object storage) buckets ACL (Access Control List).",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_acl` data source allows you to retrieve information about an existing S3 (object storage) buckets ACL (Access Control List).",
		},
		Attributes: map[string]superschema.Attribute{
			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Update: true,
					Delete: false,
					Read:   true,
				},
				DataSource: &superschema.DatasourceTimeoutAttribute{
					Read: true,
				},
			},
			"id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The ID of the S3 bucket acl. The ID is same as bucket name.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
			},
			"bucket": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Required:            true,
					MarkdownDescription: "The Name of the bucket.",
				},
				Resource: &schemaR.StringAttribute{
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
			},
			"acl": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					Optional: true,
				},
				Resource: &schemaR.StringAttribute{
					MarkdownDescription: "The canned ACL to apply to the bucket. Each canned ACL has a predefined set of grantees and permissions. [For more information](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#specifying-grantee-predefined-acl).",
					Validators: []validator.String{
						stringvalidator.OneOf("private", "public-read", "public-read-write", "authenticated-read", "bucket-owner-read", "bucket-owner-full-control"),
						stringvalidator.ConflictsWith(path.MatchRoot("access_control_policy"), path.MatchRoot("acl")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					MarkdownDescription: "The ACL applied to the bucket.",
				},
			},
			"access_control_policy": superschema.SuperSingleNestedAttributeOf[BucketACLModelAccessControlPolicy]{
				Common: &schemaR.SingleNestedAttribute{
					Computed:            true,
					MarkdownDescription: "A configuration block that sets the ACL permissions for an object per grantee.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Optional: true,
					Validators: []validator.Object{
						objectvalidator.ConflictsWith(path.MatchRoot("access_control_policy"), path.MatchRoot("acl")),
					},
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
				},
				Attributes: superschema.Attributes{
					"grants": superschema.SuperSetNestedAttributeOf[BucketACLModelGrant]{
						Common: &schemaR.SetNestedAttribute{
							MarkdownDescription: "A configuration block that sets Grant ACL permissions",
						},
						Resource: &schemaR.SetNestedAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SetNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"permission": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The following list shows each access policy permissions supported. [For more information](https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html).",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("FULL_CONTROL", "READ", "WRITE", "READ_ACP", "WRITE_ACP"),
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"grantee": superschema.SuperSingleNestedAttributeOf[BucketACLModelGrantee]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "A configuration block that sets the ACL permissions for an object per grantee.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: map[string]superschema.Attribute{
									"type": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The type of grantee specified.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf("CanonicalUser", "AmazonCustomerByEmail", "Group"),
											},
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"uri": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The URI of the grantee.",
										},
										Resource: &schemaR.StringAttribute{
											Optional: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"id": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The ID of the grantee.",
										},
										Resource: &schemaR.StringAttribute{
											Optional: true,
											// TODO - Fix validator path expression
											// Validators: []validator.String{
											// 	// fstringvalidator.RequireIfAttributeIsOneOf(path.MatchRelative().AtParent().AtName("type"), []attr.Value{types.StringValue("CanonicalUser")}),

											// },
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"display_name": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The display name of the grantee.",
											Computed:            true,
										},
									},
									"email_address": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "The email address of the grantee.",
										},
										Resource: &schemaR.StringAttribute{
											Optional: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
								},
							},
						},
					},
					"owner": superschema.SuperSingleNestedAttributeOf[BucketACLModelOwner]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "A configuration block of the bucket owner's display name and ID.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: map[string]superschema.Attribute{
							"id": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The ID of the bucket owner.",
								},
								Resource: &schemaR.StringAttribute{
									Required: true,
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"display_name": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "The display name of the bucket owner.",
									Computed:            true,
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
