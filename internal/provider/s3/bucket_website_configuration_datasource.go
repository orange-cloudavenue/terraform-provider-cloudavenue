// Package s3 provides a Terraform datasource.
package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

var (
	_ datasource.DataSource              = &BucketWebsiteConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &BucketWebsiteConfigurationDataSource{}
)

func NewBucketWebsiteConfigurationDataSource() datasource.DataSource {
	return &BucketWebsiteConfigurationDataSource{}
}

type BucketWebsiteConfigurationDataSource struct {
	client   *client.CloudAvenue
	s3Client *s3.S3
}

// Init Initializes the data source.
func (d *BucketWebsiteConfigurationDataSource) Init(ctx context.Context, dm *BucketWebsiteConfigurationDataSourceModel) (diags diag.Diagnostics) {
	d.s3Client = d.client.CAVSDK.V1.S3()
	return
}

func (d *BucketWebsiteConfigurationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_website_configuration"
}

func (d *BucketWebsiteConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bucketWebsiteConfigurationSchema(ctx).GetDataSource(ctx)
}

func (d *BucketWebsiteConfigurationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BucketWebsiteConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_s3_bucket_website_configuration", d.client.GetOrgName(), metrics.Read)()

	config := &BucketWebsiteConfigurationDataSourceModel{}

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

	data, _, diags := genericReadWebsiteConfiguration(ctx, &readWebsiteConfigurationConfig[*BucketWebsiteConfigurationDataSourceModel]{
		Client:     d.s3Client,
		BucketName: config.Bucket.GetPtr(),
	}, config)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
