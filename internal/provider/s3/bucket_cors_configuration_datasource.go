// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// Package s3 provides a Terraform datasource.
package s3

import (
	"context"
	"fmt"
	"strings"

	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &BucketCorsConfigurationDatasource{}
	_ datasource.DataSourceWithConfigure = &BucketCorsConfigurationDatasource{}
)

func NewBucketCorsConfigurationDatasource() datasource.DataSource {
	return &BucketCorsConfigurationDatasource{}
}

type BucketCorsConfigurationDatasource struct {
	client   *client.CloudAvenue
	s3Client Client
}

// Init Initializes the data source.
func (d *BucketCorsConfigurationDatasource) Init(_ context.Context, _ *BucketCorsConfigurationModelDatasource) (diags diag.Diagnostics) {
	d.s3Client, diags = initS3Client(d.client)
	return diags
}

func (d *BucketCorsConfigurationDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_cors_configuration"
}

func (d *BucketCorsConfigurationDatasource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bucketCorsConfigurationSchema(ctx).GetDataSource(ctx)
}

func (d *BucketCorsConfigurationDatasource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *BucketCorsConfigurationDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_s3_bucket_cors_configuration", d.client.GetOrgName(), metrics.Read)()

	data := &BucketCorsConfigurationModelDatasource{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set timeouts
	readTimeout, diags := data.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	/*
		Implement the data source read logic here.
	*/

	// Read the CORS
	corsResponse, err := d.s3Client.GetBucketCors(ctx, &awss3.GetBucketCorsInput{
		Bucket: data.Bucket.GetPtr(),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchCORSConfiguration):
			resp.Diagnostics.AddError("CORS policy not found", err.Error())
		default:
			resp.Diagnostics.AddError("Error retrieving CORS policy", err.Error())
		}
		return
	}
	corsPolicy := new(BucketCorsConfigurationModelCorsRules)
	for _, corsRule := range corsResponse.CORSRules {
		corsRuleModel := NewBucketCorsConfigurationModelCorsRule()

		if corsRule.MaxAgeSeconds != nil {
			corsRuleModel.MaxAgeSeconds.Set(int64(*corsRule.MaxAgeSeconds))
		}
		corsRuleModel.ID.SetPtr(corsRule.ID)

		corsRuleModel.AllowedMethods.Set(ctx, corsRule.AllowedMethods)
		corsRuleModel.AllowedOrigins.Set(ctx, corsRule.AllowedOrigins)

		// AllowedHeaders and ExposeHeaders are optional
		if len(corsRule.AllowedHeaders) > 0 {
			corsRuleModel.AllowedHeaders.Set(ctx, corsRule.AllowedHeaders)
		}
		if len(corsRule.ExposeHeaders) > 0 {
			corsRuleModel.ExposeHeaders.Set(ctx, corsRule.ExposeHeaders)
		}

		*corsPolicy = append(*corsPolicy, corsRuleModel)
	}
	resp.Diagnostics.Append(data.CorsRules.Set(ctx, corsPolicy)...)

	if !data.ID.IsKnown() {
		data.ID.Set(data.Bucket.Get())
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
