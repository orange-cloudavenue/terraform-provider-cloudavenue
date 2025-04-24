/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

// Package s3 provides a Terraform datasource.
package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &BucketLifecycleConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &BucketLifecycleConfigurationDataSource{}
)

func NewBucketLifecycleConfigurationDataSource() datasource.DataSource {
	return &BucketLifecycleConfigurationDataSource{}
}

type BucketLifecycleConfigurationDataSource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the data source.
func (d *BucketLifecycleConfigurationDataSource) Init(_ context.Context, _ *BucketLifecycleConfigurationDatasourceModel) (diags diag.Diagnostics) {
	d.s3Client = d.client.CAVSDK.V1.S3()
	return
}

func (d *BucketLifecycleConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_lifecycle_configuration"
}

func (d *BucketLifecycleConfigurationDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bucketLifecycleConfigurationSchema(ctx).GetDataSource(ctx)
}

func (d *BucketLifecycleConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BucketLifecycleConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_s3_bucket_lifecycle_configuration", d.client.GetOrgName(), metrics.Read)()

	config := &BucketLifecycleConfigurationDatasourceModel{}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(d.Init(ctx, config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the data source read logic here.
	*/

	data, _, diags := genericReadLifeCycleConfiguration(ctx, &readLifeCycleConfigurationConfig[*BucketLifecycleConfigurationDatasourceModel]{
		Client: d.s3Client.S3,
		Timeout: func() (time.Duration, diag.Diagnostics) {
			return config.Timeouts.Read(ctx, defaultReadTimeout)
		},
		BucketName: func() *string {
			return config.Bucket.GetPtr()
		},
	}, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
