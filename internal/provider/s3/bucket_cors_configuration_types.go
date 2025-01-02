package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BucketCorsConfigurationModel struct {
		Timeouts  timeoutsR.Value           `tfsdk:"timeouts"`
		ID        supertypes.StringValue    `tfsdk:"id"`
		Bucket    supertypes.StringValue    `tfsdk:"bucket"`
		CorsRules supertypes.SetNestedValue `tfsdk:"cors_rules"`
	}

	BucketCorsConfigurationModelDatasource struct {
		Timeouts  timeoutsD.Value           `tfsdk:"timeouts"`
		ID        supertypes.StringValue    `tfsdk:"id"`
		Bucket    supertypes.StringValue    `tfsdk:"bucket"`
		CorsRules supertypes.SetNestedValue `tfsdk:"cors_rules"`
	}

	BucketCorsConfigurationModelCorsRule struct {
		AllowedHeaders supertypes.SetValue    `tfsdk:"allowed_headers"`
		AllowedMethods supertypes.SetValue    `tfsdk:"allowed_methods"`
		AllowedOrigins supertypes.SetValue    `tfsdk:"allowed_origins"`
		ExposeHeaders  supertypes.SetValue    `tfsdk:"expose_headers"`
		ID             supertypes.StringValue `tfsdk:"id"`
		MaxAgeSeconds  supertypes.Int64Value  `tfsdk:"max_age_seconds"`
	}

	BucketCorsConfigurationModelCorsRules      []BucketCorsConfigurationModelCorsRule
	BucketCorsConfigurationModelAllowedHeaders []supertypes.SetValue
	BucketCorsConfigurationModelAllowedMethods []supertypes.SetValue
	BucketCorsConfigurationModelAllowedOrigins []supertypes.SetValue
	BucketCorsConfigurationModelExposeHeaders  []supertypes.SetValue
)

func (rm *BucketCorsConfigurationModel) Copy() *BucketCorsConfigurationModel {
	x := &BucketCorsConfigurationModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetCorsRules returns the value of the CorsRules field.
func (rm *BucketCorsConfigurationModel) GetCorsRules(ctx context.Context) (values BucketCorsConfigurationModelCorsRules, diags diag.Diagnostics) {
	values = make(BucketCorsConfigurationModelCorsRules, 0)
	d := rm.CorsRules.Get(ctx, &values, false)
	return values, d
}

// CorsRulesToS3CorsRules converts the BucketCorsConfigurationModelCorsRules to a slice of []*s3.CORSRule.
func (c BucketCorsConfigurationModelCorsRules) CorsRulesToS3CorsRules(ctx context.Context) (s3CorsRules []*s3.CORSRule, diags diag.Diagnostics) {
	s3CorsRules = make([]*s3.CORSRule, 0)

	for _, corsRule := range c {
		allowedHeaders := make([]string, 0)
		diags.Append(corsRule.AllowedHeaders.Get(ctx, &allowedHeaders, false)...)

		allowedMethods := make([]string, 0)
		diags.Append(corsRule.AllowedMethods.Get(ctx, &allowedMethods, false)...)

		allowedOrigins := make([]string, 0)
		diags.Append(corsRule.AllowedOrigins.Get(ctx, &allowedOrigins, false)...)

		exposeHeaders := make([]string, 0)
		diags.Append(corsRule.ExposeHeaders.Get(ctx, &exposeHeaders, false)...)

		if diags.HasError() {
			return nil, diags
		}

		s3CorsRule := &s3.CORSRule{
			ID:             corsRule.ID.GetPtr(),
			MaxAgeSeconds:  corsRule.MaxAgeSeconds.GetPtr(),
			AllowedHeaders: utils.SliceToSlicePointer(allowedHeaders),
			AllowedMethods: utils.SliceToSlicePointer(allowedMethods),
			AllowedOrigins: utils.SliceToSlicePointer(allowedOrigins),
			ExposeHeaders:  utils.SliceToSlicePointer(exposeHeaders),
		}

		s3CorsRules = append(s3CorsRules, s3CorsRule)
	}

	return s3CorsRules, nil
}

// NewBucketCorsConfigurationModelCorsRule creates a new BucketCorsConfigurationModelCorsRule.
func NewBucketCorsConfigurationModelCorsRule() BucketCorsConfigurationModelCorsRule {
	return BucketCorsConfigurationModelCorsRule{
		AllowedHeaders: supertypes.NewSetNull(supertypes.StringType{}),
		AllowedMethods: supertypes.NewSetNull(supertypes.StringType{}),
		AllowedOrigins: supertypes.NewSetNull(supertypes.StringType{}),
		ExposeHeaders:  supertypes.NewSetNull(supertypes.StringType{}),
		ID:             supertypes.NewStringNull(),
		MaxAgeSeconds:  supertypes.NewInt64Null(),
	}
}
