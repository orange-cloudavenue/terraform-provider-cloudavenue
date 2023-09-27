package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

func providerSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		MarkdownDescription: "This provider offers utilities for working with the Cloud Avenue platform.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The URL of the Cloud Avenue API. Can also be set with the `CLOUDAVENUE_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https?:\/\/\S+\w$`),
						"must end with a letter",
					),
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "The username to use to connect to the Cloud Avenue API. Can also be set with the `CLOUDAVENUE_USER` environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to use to connect to the Cloud Avenue API. Can also be set with the `CLOUDAVENUE_PASSWORD` environment variable.",
				Sensitive:           true,
				Optional:            true,
			},
			"org": schema.StringAttribute{
				MarkdownDescription: "The organization used on Cloud Avenue API. Can also be set with the `CLOUDAVENUE_ORG` environment variable.",
				Optional:            true,
			},
			"vdc": schema.StringAttribute{
				MarkdownDescription: "The VDC used on Cloud Avenue API. Can also be set with the `CLOUDAVENUE_VDC` environment variable.",
				Optional:            true,
			},
			"netbackup_url": schema.StringAttribute{
				MarkdownDescription: "The URL of the NetBackup API. Can also be set with the `NETBACKUP_URL` environment variable.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https?:\/\/\S+\w$`),
						"must end with a letter",
					),
				},
			},
			"netbackup_user": schema.StringAttribute{
				MarkdownDescription: "The username to use to connect to the NetBackup API. Can also be set with the `NETBACKUP_USER` environment variable.",
				Optional:            true,
			},
			"netbackup_password": schema.StringAttribute{
				MarkdownDescription: "The password to use to connect to the NetBackup API. Can also be set with the `NETBACKUP_PASSWORD` environment variable.",
				Sensitive:           true,
				Optional:            true,
			},
		},
	}
}
