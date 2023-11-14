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
	_ datasource.DataSource              = &BucketPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &BucketPolicyDataSource{}
)

func NewBucketPolicyDataSource() datasource.DataSource {
	return &BucketPolicyDataSource{}
}

type BucketPolicyDataSource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the data source.
func (d *BucketPolicyDataSource) Init(ctx context.Context, dm *BucketPolicyModelDatasource) (diags diag.Diagnostics) {
	d.s3Client = d.client.CAVSDK.V1.S3()
	return
}

func (d *BucketPolicyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_policy"
}

func (d *BucketPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bucketPolicySchema(ctx).GetDataSource(ctx)
}

func (d *BucketPolicyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *BucketPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer metrics.New("data.cloudavenue_s3_bucket_policy", d.client.GetOrgName(), metrics.Read)()

	config := &BucketPolicyModelDatasource{}

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

	// Set timeouts
	readTimeout, diags := config.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	// Get Policy
	policy, err := findBucketPolicy(ctx, &readPolicyConfig[*BucketPolicyModel]{
		Timeout: func() (time.Duration, diag.Diagnostics) {
			return readTimeout, diags
		},
		Client: d.s3Client.S3,
		BucketName: func() *string {
			return config.Bucket.GetPtr()
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Bucket Policy not found", err.Error())
		return
	}

	config.ID.Set(config.Bucket.Get())
	config.Policy.Set(policy)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
