package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BucketWebsiteConfigurationModel struct {
		Timeouts              timeoutsR.Value                                                                            `tfsdk:"timeouts"`
		ID                    supertypes.StringValue                                                                     `tfsdk:"id"`
		Bucket                supertypes.StringValue                                                                     `tfsdk:"bucket"`
		ErrorDocument         supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelErrorDocument]         `tfsdk:"error_document"`
		IndexDocument         supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelIndexDocument]         `tfsdk:"index_document"`
		RedirectAllRequestsTo supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelRedirectAllRequestsTo] `tfsdk:"redirect_all_requests_to"`
		RoutingRules          supertypes.SetNestedObjectValueOf[BucketWebsiteConfigurationModelRoutingRule]              `tfsdk:"routing_rules"`
		WebsiteEndpoint       supertypes.StringValue                                                                     `tfsdk:"website_endpoint"`
	}

	BucketWebsiteConfigurationDataSourceModel struct {
		Timeouts              timeoutsD.Value                                                                            `tfsdk:"timeouts"`
		ID                    supertypes.StringValue                                                                     `tfsdk:"id"`
		Bucket                supertypes.StringValue                                                                     `tfsdk:"bucket"`
		ErrorDocument         supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelErrorDocument]         `tfsdk:"error_document"`
		IndexDocument         supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelIndexDocument]         `tfsdk:"index_document"`
		RedirectAllRequestsTo supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelRedirectAllRequestsTo] `tfsdk:"redirect_all_requests_to"`
		RoutingRules          supertypes.SetNestedObjectValueOf[BucketWebsiteConfigurationModelRoutingRule]              `tfsdk:"routing_rules"`
		WebsiteEndpoint       supertypes.StringValue                                                                     `tfsdk:"website_endpoint"`
	}

	BucketWebsiteConfigurationModelErrorDocument struct {
		Key supertypes.StringValue `tfsdk:"key"`
	}

	BucketWebsiteConfigurationModelIndexDocument struct {
		Suffix supertypes.StringValue `tfsdk:"suffix"`
	}

	BucketWebsiteConfigurationModelRedirectAllRequestsTo struct {
		HostName supertypes.StringValue `tfsdk:"hostname"`
		Protocol supertypes.StringValue `tfsdk:"protocol"`
	}

	BucketWebsiteConfigurationModelRoutingRule struct {
		Condition supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelCondition] `tfsdk:"condition"`
		Redirect  supertypes.SingleNestedObjectValueOf[BucketWebsiteConfigurationModelRedirect]  `tfsdk:"redirect"`
	}

	BucketWebsiteConfigurationModelCondition struct {
		HTTPErrorCodeReturnedEquals supertypes.StringValue `tfsdk:"http_error_code_returned_equals"`
		KeyPrefixEquals             supertypes.StringValue `tfsdk:"key_prefix_equals"`
	}

	BucketWebsiteConfigurationModelRedirect struct {
		HostName             supertypes.StringValue `tfsdk:"hostname"`
		HTTPRedirectCode     supertypes.StringValue `tfsdk:"http_redirect_code"`
		Protocol             supertypes.StringValue `tfsdk:"protocol"`
		ReplaceKeyPrefixWith supertypes.StringValue `tfsdk:"replace_key_prefix_with"`
		ReplaceKeyWith       supertypes.StringValue `tfsdk:"replace_key_with"`
	}
)

func (rm *BucketWebsiteConfigurationModel) CreateS3WebsiteConfigurationAPIObject(ctx context.Context) (websiteConfig *s3.WebsiteConfiguration, diags diag.Diagnostics) {
	websiteConfig = new(s3.WebsiteConfiguration)

	if rm.ErrorDocument.IsKnown() {
		errorDocument, d := rm.ErrorDocument.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
		websiteConfig.ErrorDocument = &s3.ErrorDocument{
			Key: errorDocument.Key.GetPtr(),
		}
	}

	if rm.IndexDocument.IsKnown() {
		indexDocument, d := rm.IndexDocument.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		websiteConfig.IndexDocument = &s3.IndexDocument{
			Suffix: indexDocument.Suffix.GetPtr(),
		}
	}

	if rm.RedirectAllRequestsTo.IsKnown() {
		redirectAllRequestsTo, d := rm.RedirectAllRequestsTo.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		websiteConfig.RedirectAllRequestsTo = &s3.RedirectAllRequestsTo{
			HostName: redirectAllRequestsTo.HostName.GetPtr(),
			Protocol: redirectAllRequestsTo.Protocol.GetPtr(),
		}
	}

	if rm.RoutingRules.IsKnown() {
		routingRules, d := rm.RoutingRules.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return
		}

		websiteConfig.RoutingRules = make([]*s3.RoutingRule, len(routingRules))
		for i, routingRule := range routingRules {
			redirect, d := routingRule.Redirect.Get(ctx)
			diags.Append(d...)
			condition, d := routingRule.Condition.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return
			}

			websiteConfig.RoutingRules[i] = &s3.RoutingRule{
				Redirect: &s3.Redirect{
					HostName:             redirect.HostName.GetPtr(),
					HttpRedirectCode:     redirect.HTTPRedirectCode.GetPtr(),
					Protocol:             redirect.Protocol.GetPtr(),
					ReplaceKeyPrefixWith: redirect.ReplaceKeyPrefixWith.GetPtr(),
					ReplaceKeyWith:       redirect.ReplaceKeyWith.GetPtr(),
				},
				Condition: &s3.Condition{
					HttpErrorCodeReturnedEquals: condition.HTTPErrorCodeReturnedEquals.GetPtr(),
					KeyPrefixEquals:             condition.KeyPrefixEquals.GetPtr(),
				},
			}
		}
	}

	return
}

func (rm *BucketWebsiteConfigurationModel) SetErrorDocument(ctx context.Context, errorDocument *BucketWebsiteConfigurationModelErrorDocument) diag.Diagnostics {
	return rm.ErrorDocument.Set(ctx, errorDocument)
}

func (dm *BucketWebsiteConfigurationDataSourceModel) SetErrorDocument(ctx context.Context, errorDocument *BucketWebsiteConfigurationModelErrorDocument) diag.Diagnostics {
	return dm.ErrorDocument.Set(ctx, errorDocument)
}

func (rm *BucketWebsiteConfigurationModel) SetIndexDocument(ctx context.Context, indexDocument *BucketWebsiteConfigurationModelIndexDocument) diag.Diagnostics {
	return rm.IndexDocument.Set(ctx, indexDocument)
}

func (dm *BucketWebsiteConfigurationDataSourceModel) SetIndexDocument(ctx context.Context, indexDocument *BucketWebsiteConfigurationModelIndexDocument) diag.Diagnostics {
	return dm.IndexDocument.Set(ctx, indexDocument)
}

func (rm *BucketWebsiteConfigurationModel) SetRedirectAllRequestsTo(ctx context.Context, redirectAllRequestsTo *BucketWebsiteConfigurationModelRedirectAllRequestsTo) diag.Diagnostics {
	return rm.RedirectAllRequestsTo.Set(ctx, redirectAllRequestsTo)
}

func (dm *BucketWebsiteConfigurationDataSourceModel) SetRedirectAllRequestsTo(ctx context.Context, redirectAllRequestsTo *BucketWebsiteConfigurationModelRedirectAllRequestsTo) diag.Diagnostics {
	return dm.RedirectAllRequestsTo.Set(ctx, redirectAllRequestsTo)
}

func (rm *BucketWebsiteConfigurationModel) SetRoutingRules(ctx context.Context, routingRules []*BucketWebsiteConfigurationModelRoutingRule) diag.Diagnostics {
	return rm.RoutingRules.Set(ctx, routingRules)
}

func (dm *BucketWebsiteConfigurationDataSourceModel) SetRoutingRules(ctx context.Context, routingRules []*BucketWebsiteConfigurationModelRoutingRule) diag.Diagnostics {
	return dm.RoutingRules.Set(ctx, routingRules)
}

func (rm *BucketWebsiteConfigurationModel) SetID(id string) {
	rm.ID.Set(id)
}

func (dm *BucketWebsiteConfigurationDataSourceModel) SetID(id string) {
	dm.ID.Set(id)
}

func (rm *BucketWebsiteConfigurationModel) Copy() any {
	x := &BucketWebsiteConfigurationModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (dm *BucketWebsiteConfigurationDataSourceModel) Copy() any {
	x := &BucketWebsiteConfigurationDataSourceModel{}
	utils.ModelCopy(dm, x)
	return x
}

func (rm *BucketWebsiteConfigurationModel) SetWebsiteEnpoint(endpoint string) {
	rm.WebsiteEndpoint.Set(endpoint)
}

func (dm *BucketWebsiteConfigurationDataSourceModel) SetWebsiteEnpoint(endpoint string) {
	dm.WebsiteEndpoint.Set(endpoint)
}

type (
	readWebsiteConfigurationResourceDatasource interface {
		*BucketWebsiteConfigurationDataSourceModel | *BucketWebsiteConfigurationModel
		SetID(id string)
		SetErrorDocument(ctx context.Context, errorDocument *BucketWebsiteConfigurationModelErrorDocument) diag.Diagnostics
		SetIndexDocument(ctx context.Context, indexDocument *BucketWebsiteConfigurationModelIndexDocument) diag.Diagnostics
		SetRedirectAllRequestsTo(ctx context.Context, redirectAllRequestsTo *BucketWebsiteConfigurationModelRedirectAllRequestsTo) diag.Diagnostics
		SetRoutingRules(ctx context.Context, routingRules []*BucketWebsiteConfigurationModelRoutingRule) diag.Diagnostics
		Copy() any
		SetWebsiteEnpoint(endpoint string)
	}

	readWebsiteConfigurationConfig[T readWebsiteConfigurationResourceDatasource] struct {
		Client     *s3.S3
		BucketName *string
	}
)

// genericReadWebsiteConfiguration is a generic function that reads the website configuration of a bucket.
func genericReadWebsiteConfiguration[T readWebsiteConfigurationResourceDatasource](ctx context.Context, config *readWebsiteConfigurationConfig[T], planOrState T) (stateRefreshed T, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy().(T)

	output, err := findBucketWebsite(ctx, config.Client, *config.BucketName)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchBucket):
			diags.AddError("Bucket not found", err.Error())
			return nil, false, diags
		default:
			diags.AddError("Error getting website configuration", err.Error())
			return nil, true, diags
		}
	}

	if output == nil {
		diags.AddError("Unexpected nil output from GetBucketWebsite", "")
		return nil, true, diags
	}

	if output.ErrorDocument != nil {
		eD := &BucketWebsiteConfigurationModelErrorDocument{
			Key: supertypes.NewStringNull(),
		}
		eD.Key.SetPtr(output.ErrorDocument.Key)
		diags.Append(stateRefreshed.SetErrorDocument(ctx, eD)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	if output.IndexDocument != nil {
		iD := &BucketWebsiteConfigurationModelIndexDocument{
			Suffix: supertypes.NewStringNull(),
		}

		iD.Suffix.SetPtr(output.IndexDocument.Suffix)
		diags.Append(stateRefreshed.SetIndexDocument(ctx, iD)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	if output.RedirectAllRequestsTo != nil {
		rART := &BucketWebsiteConfigurationModelRedirectAllRequestsTo{
			HostName: supertypes.NewStringNull(),
			Protocol: supertypes.NewStringNull(),
		}

		rART.HostName.SetPtr(output.RedirectAllRequestsTo.HostName)
		rART.Protocol.SetPtr(output.RedirectAllRequestsTo.Protocol)
		diags.Append(stateRefreshed.SetRedirectAllRequestsTo(ctx, rART)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	if output.RoutingRules != nil {
		rRs := make([]*BucketWebsiteConfigurationModelRoutingRule, 0, len(output.RoutingRules))

		for _, rule := range output.RoutingRules {
			redirect := &BucketWebsiteConfigurationModelRedirect{
				HostName:             supertypes.NewStringNull(),
				HTTPRedirectCode:     supertypes.NewStringNull(),
				Protocol:             supertypes.NewStringNull(),
				ReplaceKeyPrefixWith: supertypes.NewStringNull(),
				ReplaceKeyWith:       supertypes.NewStringNull(),
			}

			condition := &BucketWebsiteConfigurationModelCondition{
				HTTPErrorCodeReturnedEquals: supertypes.NewStringNull(),
				KeyPrefixEquals:             supertypes.NewStringNull(),
			}

			rR := &BucketWebsiteConfigurationModelRoutingRule{
				Redirect:  supertypes.NewSingleNestedObjectValueOfNull[BucketWebsiteConfigurationModelRedirect](ctx),
				Condition: supertypes.NewSingleNestedObjectValueOfNull[BucketWebsiteConfigurationModelCondition](ctx),
			}

			if rule.Redirect != nil {
				redirect.HostName.SetPtr(rule.Redirect.HostName)
				redirect.HTTPRedirectCode.SetPtr(rule.Redirect.HttpRedirectCode)
				redirect.Protocol.SetPtr(rule.Redirect.Protocol)
				redirect.ReplaceKeyPrefixWith.SetPtr(rule.Redirect.ReplaceKeyPrefixWith)
				redirect.ReplaceKeyWith.SetPtr(rule.Redirect.ReplaceKeyWith)
				diags.Append(rR.Redirect.Set(ctx, redirect)...)
				if diags.HasError() {
					return nil, true, diags
				}
			}

			if rule.Condition != nil {
				condition.HTTPErrorCodeReturnedEquals.SetPtr(rule.Condition.HttpErrorCodeReturnedEquals)
				condition.KeyPrefixEquals.SetPtr(rule.Condition.KeyPrefixEquals)
				diags.Append(rR.Condition.Set(ctx, condition)...)
				if diags.HasError() {
					return nil, true, diags
				}
			}

			rRs = append(rRs, rR)
		}

		diags.Append(stateRefreshed.SetRoutingRules(ctx, rRs)...)
		if diags.HasError() {
			return nil, true, diags
		}
	}

	stateRefreshed.SetID(*config.BucketName)
	stateRefreshed.SetWebsiteEnpoint(fmt.Sprintf("%s.%s", *config.BucketName, "website-region01.cloudavenue.orange-business.com"))

	return stateRefreshed, true, nil
}
