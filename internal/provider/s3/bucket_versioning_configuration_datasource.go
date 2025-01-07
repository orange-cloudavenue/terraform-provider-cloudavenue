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

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &BucketVersioningConfigurationDatasource{}
	_ datasource.DataSourceWithConfigure = &BucketVersioningConfigurationDatasource{}
)

func NewBucketVersioningConfigurationDatasource() datasource.DataSource {
	return &BucketVersioningConfigurationDatasource{}
}

type BucketVersioningConfigurationDatasource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the data source.
func (d *BucketVersioningConfigurationDatasource) Init(ctx context.Context, dm *BucketVersioningConfigurationDatasourceModel) (diags diag.Diagnostics) {
	d.s3Client = d.client.CAVSDK.V1.S3()
	return
}

func (d *BucketVersioningConfigurationDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_versioning_configuration"
}

func (d *BucketVersioningConfigurationDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bucketVersioningConfigurationSchema(ctx).GetDataSource(ctx)
}

func (d *BucketVersioningConfigurationDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BucketVersioningConfigurationDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_s3_bucket_versioning_configuration", d.client.GetOrgName(), metrics.Read)()

	// If the data source don't have same schema/structure as the resource, you can use the following code:
	data := &BucketVersioningConfigurationDatasourceModel{}

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

	/*
		Implement the data source read logic here.
	*/

	// Read the versioning configuration
	versioningResponse, err := d.s3Client.GetBucketVersioningWithContext(ctx, &s3.GetBucketVersioningInput{
		Bucket: data.Bucket.GetPtr(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving bucket versioning configuration", err.Error())
		return
	}

	data.Status.SetPtr(versioningResponse.Status)

	if !data.ID.IsKnown() {
		data.ID.Set(data.Bucket.Get())
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
