package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketPolicyResource{}
	_ resource.ResourceWithConfigure   = &BucketPolicyResource{}
	_ resource.ResourceWithImportState = &BucketPolicyResource{}
)

// NewBucketPolicyResource is a helper function to simplify the provider implementation.
func NewBucketPolicyResource() resource.Resource {
	return &BucketPolicyResource{}
}

// BucketPolicyResource is the resource implementation.
type BucketPolicyResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketPolicyResource) Init(ctx context.Context, rm *BucketPolicyModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()

	return
}

// Metadata returns the resource type name.
func (r *BucketPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_policy"
}

// Schema defines the schema for the resource.
func (r *BucketPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketPolicySchema(ctx).GetResource(ctx)
}

func (r *BucketPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.CloudAvenue)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.CloudAvenue, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *BucketPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_policy", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketPolicyModel{}

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set timeouts
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	/*
		Implement the resource creation logic here.
	*/

	stateRefreshed, d := r.genericCreateOrUpdate(ctx, createTimeout, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_policy", r.client.GetOrgName(), metrics.Read)()

	state := &BucketPolicyModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set timeouts
	readTimeout, diags := state.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	// Refresh the state
	stateRefreshed, found, diags := r.read(ctx, readTimeout, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	_, err := policyToSet(state.GetPolicy(ctx), stateRefreshed.GetPolicy(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Bucket Policy JSON error", err.Error())
		return
	}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BucketPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update for this resource, only delete and create for a change.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_policy", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketPolicyModel{}

	// Get current state
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	/*
		Implement the resource deletion here
	*/

	// Set timeouts
	deleteTimeout, d := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	input := &s3.DeleteBucketPolicyInput{
		Bucket: state.Bucket.GetPtr(),
	}

	// Delete Policy
	_, err := retryWhenAWSErrCodeEquals(ctx, []string{ErrCodeNoSuchBucket}, &RetryWhenConfig[*s3.DeleteBucketPolicyOutput]{
		Timeout:  deleteTimeout,
		Interval: 15 * time.Second,
		Function: func() (*s3.DeleteBucketPolicyOutput, error) {
			return r.s3Client.DeleteBucketPolicy(input)
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Deleting S3 Bucket Policy", err.Error())
		return
	}
}

func (r *BucketPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_policy", r.client.GetOrgName(), metrics.Import)()

	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// genericCreateOrUpdate creates or updates a resource.
func (r *BucketPolicyResource) genericCreateOrUpdate(ctx context.Context, timeout time.Duration, planOrState *BucketPolicyModel) (stateRefreshed *BucketPolicyModel, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy().(*BucketPolicyModel)

	// Check JSON Policy
	policy, err := structure.NormalizeJsonString(stateRefreshed.Policy.Get())
	if err != nil {
		diags.AddError("Invalid policy", err.Error())
		return nil, diags
	}

	// Loop for object creation
	input := &s3.PutBucketPolicyInput{
		Bucket: stateRefreshed.Bucket.GetPtr(),
		Policy: aws.String(policy),
	}
	_, err = retryWhenAWSErrCodeEquals(ctx, []string{ErrCodeNoSuchBucket, ErrCodeMalformedPolicy}, &RetryWhenConfig[*s3.PutBucketPolicyOutput]{
		Timeout:  timeout,
		Interval: 15 * time.Second,
		Function: func() (*s3.PutBucketPolicyOutput, error) {
			return r.s3Client.PutBucketPolicy(input)
		},
	})
	if err != nil {
		diags.AddError("Bucket Policy Creation Error", err.Error())
		return nil, diags
	}

	stateRefreshed, _, diags = r.read(ctx, timeout, stateRefreshed)

	return stateRefreshed, diags
}

// read function for bucket policy.
func (r *BucketPolicyResource) read(ctx context.Context, timeout time.Duration, planOrState *BucketPolicyModel) (stateRefreshed *BucketPolicyModel, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy().(*BucketPolicyModel)

	// Get Policy
	policy, err := findBucketPolicy(ctx, &readPolicyConfig[*BucketPolicyModel]{
		Timeout: func() (time.Duration, diag.Diagnostics) {
			return timeout, diags
		},
		Client: r.s3Client.S3,
		BucketName: func() *string {
			return stateRefreshed.Bucket.GetPtr()
		},
	})
	if err != nil {
		diags.AddError("Bucket Policy not found", err.Error())
		return nil, false, diags
	}

	// Set the ID and Policy.
	stateRefreshed.SetID(stateRefreshed.Bucket.GetPtr())
	stateRefreshed.SetPolicy(policy)

	return stateRefreshed, true, diags
}
