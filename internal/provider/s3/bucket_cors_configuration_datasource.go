// Package s3 provides a Terraform datasource.
package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
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
	s3Client v1.S3Client
}

// Init Initializes the data source.
func (d *BucketCorsConfigurationDatasource) Init(ctx context.Context, dm *BucketCorsConfigurationModelDatasource) (diags diag.Diagnostics) {
	d.s3Client = d.client.CAVSDK.V1.S3()
	return
}

func (d *BucketCorsConfigurationDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_cors_configuration"
}

func (d *BucketCorsConfigurationDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = bucketCorsConfigurationSchema(ctx).GetDataSource(ctx)
}

func (d *BucketCorsConfigurationDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	corsResponse, err := d.s3Client.GetBucketCorsWithContext(ctx, &s3.GetBucketCorsInput{
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

		corsRuleModel.MaxAgeSeconds.SetPtr(corsRule.MaxAgeSeconds)
		corsRuleModel.ID.SetPtr(corsRule.ID)

		corsRuleModel.AllowedMethods.Set(ctx, utils.SlicePointerToSlice(corsRule.AllowedMethods))
		corsRuleModel.AllowedOrigins.Set(ctx, utils.SlicePointerToSlice(corsRule.AllowedOrigins))

		// AllowedHeaders and ExposeHeaders are optional
		if len(corsRule.AllowedHeaders) > 0 {
			corsRuleModel.AllowedHeaders.Set(ctx, utils.SlicePointerToSlice(corsRule.AllowedHeaders))
		}
		if len(corsRule.ExposeHeaders) > 0 {
			corsRuleModel.ExposeHeaders.Set(ctx, utils.SlicePointerToSlice(corsRule.ExposeHeaders))
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
