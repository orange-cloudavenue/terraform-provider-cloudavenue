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
	_ datasource.DataSource              = &BucketOwnershipControlsDataSource{}
	_ datasource.DataSourceWithConfigure = &BucketOwnershipControlsDataSource{}
)

func NewBucketOwnershipControlsDataSource() datasource.DataSource {
	return &BucketOwnershipControlsDataSource{}
}

type BucketOwnershipControlsDataSource struct {
	client   *client.CloudAvenue
	s3Client *s3.S3
}

// Init Initializes the data source.
func (d *BucketOwnershipControlsDataSource) Init(ctx context.Context, dm *BucketOwnershipControlsDataSourceModel) (diags diag.Diagnostics) {
	d.s3Client = d.client.CAVSDK.V1.S3()
	return
}

func (d *BucketOwnershipControlsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_ownership_controls"
}

func (d *BucketOwnershipControlsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = BucketOwnershipControlsSchema(ctx).GetDataSource(ctx)
}

func (d *BucketOwnershipControlsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BucketOwnershipControlsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_s3_bucket_ownership_controls", d.client.GetOrgName(), metrics.Read)()

	config := &BucketOwnershipControlsDataSourceModel{}

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

	stateRefreshed, found, diags := genericReadOwnerShipControls(ctx, &genericOwnershipControlsConfig[*BucketOwnershipControlsDataSourceModel]{
		Client:     d.s3Client,
		BucketName: func() *string { return config.Bucket.GetPtr() },
	}, config)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateRefreshed)...)
}
