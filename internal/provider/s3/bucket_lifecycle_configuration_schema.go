/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package s3

import (
	"context"

	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func bucketLifecycleConfigurationSchema(_ context.Context) superschema.Schema {
	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_lifecycle_configuration` resource allows you to manage lifecycle configuration of an S3 bucket.",
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: "The `cloudavenue_s3_bucket_lifecycle_configuration` data source allows you to retrieve information about an S3 bucket's lifecycle configuration.",
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
					MarkdownDescription: "The ID is a bucket name.",
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
			"rules": superschema.SuperListNestedAttributeOf[BucketLifecycleConfigurationModelRule]{
				Common: &schemaR.ListNestedAttribute{
					MarkdownDescription: "Rules that define lifecycle configuration.",
				},
				Resource: &schemaR.ListNestedAttribute{
					Required: true,
					Validators: []validator.List{
						listvalidator.SizeAtLeast(1),
					},
				},
				DataSource: &schemaD.ListNestedAttribute{
					Computed: true,
				},
				Attributes: superschema.Attributes{
					"abort_incomplete_multipart_upload": superschema.SuperSingleNestedAttributeOf[BucketLifecycleConfigurationModelAbortIncompleteMultipartUpload]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Configuration block that specifies the days since the initiation of an incomplete multipart upload that S3 will wait before permanently removing all parts of the upload",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"days_after_initiation": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "Number of days after which S3 aborts an incomplete multipart upload.",
									Optional:            true,
								},
							},
						},
					},
					"expiration": superschema.SuperSingleNestedAttributeOf[BucketLifecycleConfigurationModelExpiration]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Configuration block that specifies the expiration for the lifecycle of the object in the form of date, days and, whether the object has a delete marker.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("expiration"), path.MatchRelative().AtParent().AtName("noncurrent_version_expiration")),
							},
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"date": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Date the object is to be moved or deleted. The date value must be in [RFC3339 full-date format](https://datatracker.ietf.org/doc/html/rfc3339#section-5.6) e.g. `2023-10-10T00:00:00Z`",
								},
								Resource: &schemaR.StringAttribute{
									Optional: true,
									Validators: []validator.String{
										stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("date"), path.MatchRelative().AtParent().AtName("days"), path.MatchRelative().AtParent().AtName("expired_object_delete_marker")),
										// TODO : Add validator for RFC3339 full-date format.
									},
								},
								DataSource: &schemaD.StringAttribute{
									Computed: true,
								},
							},
							"days": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "Lifetime, in days, of the objects that are subject to the rule. The value must be a non-zero positive integer.",
									Computed:            true, // API return 0 if not set.
								},
								Resource: &schemaR.Int64Attribute{
									Optional: true,
									Validators: []validator.Int64{
										int64validator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("date"), path.MatchRelative().AtParent().AtName("days"), path.MatchRelative().AtParent().AtName("expired_object_delete_marker")),
									},
								},
							},
							"expired_object_delete_marker": superschema.SuperBoolAttribute{
								Common: &schemaR.BoolAttribute{
									MarkdownDescription: "Indicates whether S3 will remove a delete marker with no noncurrent versions. If set to `true`, the delete marker will be expired, if set to `false` the policy takes no action",
									Computed:            true, // API return false if not set.
								},
								Resource: &schemaR.BoolAttribute{
									Optional: true,
									Validators: []validator.Bool{
										boolvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("date"), path.MatchRelative().AtParent().AtName("days"), path.MatchRelative().AtParent().AtName("expired_object_delete_marker")),
										boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtParent().AtName("filter").AtName("tag"), path.MatchRelative().AtParent().AtParent().AtName("filter").AtName("and").AtName("tags")),
									},
								},
							},
						},
					},
					"filter": superschema.SuperSingleNestedAttributeOf[BucketLifecycleConfigurationModelFilter]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Configuration block used to identify objects that a Lifecycle Rule applies to.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Required: true,
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"and": superschema.SuperSingleNestedAttributeOf[BucketLifecycleConfigurationModelAnd]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Configuration block used to apply a logical AND to two or more predicates. The Lifecycle Rule will apply to any object matching all the predicates configured inside the and block.",
								},
								Resource: &schemaR.SingleNestedAttribute{
									Optional: true,
								},
								DataSource: &schemaD.SingleNestedAttribute{
									Computed: true,
								},
								Attributes: superschema.Attributes{
									"tags": superschema.SuperListNestedAttributeOf[BucketLifecycleConfigurationModelTag]{
										Common: &schemaR.ListNestedAttribute{
											MarkdownDescription: "Specifies object tag key and value.",
										},
										Resource: &schemaR.ListNestedAttribute{
											Optional: true,
										},
										DataSource: &schemaD.ListNestedAttribute{
											Computed: true,
										},
										Attributes: superschema.Attributes{
											"key": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "Object tag key.",
												},
												Resource: &schemaR.StringAttribute{
													Required: true,
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
											"value": superschema.SuperStringAttribute{
												Common: &schemaR.StringAttribute{
													MarkdownDescription: "Object tag value.",
												},
												Resource: &schemaR.StringAttribute{
													Required: true,
												},
												DataSource: &schemaD.StringAttribute{
													Computed: true,
												},
											},
										},
									},
									"prefix": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Match objects with this prefix.",
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
							"tag": superschema.SuperSingleNestedAttributeOf[BucketLifecycleConfigurationModelTag]{
								Common: &schemaR.SingleNestedAttribute{
									MarkdownDescription: "Specifies object tag key and value.",
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
											MarkdownDescription: "Object tag key.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
									"value": superschema.SuperStringAttribute{
										Common: &schemaR.StringAttribute{
											MarkdownDescription: "Object tag value.",
										},
										Resource: &schemaR.StringAttribute{
											Required: true,
										},
										DataSource: &schemaD.StringAttribute{
											Computed: true,
										},
									},
								},
							},
							"prefix": superschema.SuperStringAttribute{
								Common: &schemaR.StringAttribute{
									MarkdownDescription: "Match objects with this prefix.",
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
					"id": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Unique identifier for the rule.",
						},
						Resource: &schemaR.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 255),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					// TODO OneOf Expiration, NoncurrentVersionExpiration
					"noncurrent_version_expiration": superschema.SuperSingleNestedAttributeOf[BucketLifecycleConfigurationModelNoncurrentVersionExpiration]{
						Common: &schemaR.SingleNestedAttribute{
							MarkdownDescription: "Configuration block that specifies when noncurrent object versions expire.",
						},
						Resource: &schemaR.SingleNestedAttribute{
							Optional: true,
							Validators: []validator.Object{
								objectvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("expiration"), path.MatchRelative().AtParent().AtName("noncurrent_version_expiration")),
							},
						},
						DataSource: &schemaD.SingleNestedAttribute{
							Computed: true,
						},
						Attributes: superschema.Attributes{
							"noncurrent_days": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "Number of days an object is noncurrent before S3 can perform the associated action. Must be a positive integer.",
								},
								Resource: &schemaR.Int64Attribute{
									Optional: true,
									Validators: []validator.Int64{
										int64validator.AtLeast(1),
									},
								},
								DataSource: &schemaD.Int64Attribute{
									Computed: true,
								},
							},
							"newer_noncurrent_versions": superschema.SuperInt64Attribute{
								Common: &schemaR.Int64Attribute{
									MarkdownDescription: "Number of noncurrent versions S3 will retain.",
								},
								Resource: &schemaR.Int64Attribute{
									Optional: true,
									Validators: []validator.Int64{
										int64validator.AtLeast(0),
									},
								},
								DataSource: &schemaD.Int64Attribute{
									Computed: true,
								},
							},
						},
					},
					"status": superschema.SuperStringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "Whether the rule is currently being applied",
							Computed:            true,
						},
						Resource: &schemaR.StringAttribute{
							Optional: true,
							Default:  stringdefault.StaticString(lifeCycleStatusEnabled),
							Validators: []validator.String{
								stringvalidator.OneOf(lifeCycleStatusEnabled, lifeCycleStatusDisabled),
							},
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}
