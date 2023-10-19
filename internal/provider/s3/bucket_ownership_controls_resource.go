package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/k0kubun/pp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketOwnershipControlsResource{}
	_ resource.ResourceWithConfigure   = &BucketOwnershipControlsResource{}
	_ resource.ResourceWithImportState = &BucketOwnershipControlsResource{}
)

// NewBucketOwnershipControlsResource is a helper function to simplify the provider implementation.
func NewBucketOwnershipControlsResource() resource.Resource {
	return &BucketOwnershipControlsResource{}
}

// BucketOwnershipControlsResource is the resource implementation.
type BucketOwnershipControlsResource struct {
	client   *client.CloudAvenue
	s3Client *s3.S3
}

// Init Initializes the resource.
func (r *BucketOwnershipControlsResource) Init(ctx context.Context, rm *BucketOwnershipControlsModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return
}

// Metadata returns the resource type name.
func (r *BucketOwnershipControlsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_ownership_controls"
}

// Schema defines the schema for the resource.
func (r *BucketOwnershipControlsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = BucketOwnershipControlsSchema(ctx).GetResource(ctx)
}

func (r *BucketOwnershipControlsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketOwnershipControlsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_ownership_controls", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketOwnershipControlsModel{}

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

	// Set default create timeout
	createTimeout, diags := plan.Timeouts.Create(ctx, defaultCreateTimeout)
	diags.Append(diags...)
	if diags.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	/*
		Implement the resource creation logic here.
	*/

	// Create the bucket ownership controls
	d := r.createOrUpdate(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	tflog.Debug(ctx, pp.Sprintf("------Plan After func create:", plan))

	// Use generic read function to refresh the state
	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	tflog.Debug(ctx, pp.Sprintf("------State After func read:", state))

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketOwnershipControlsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_ownership_controls", r.client.GetOrgName(), metrics.Read)()

	state := &BucketOwnershipControlsModel{}

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

	// Set default read timeout
	readTimeout, diags := state.Timeouts.Read(ctx, defaultReadTimeout)
	diags.Append(diags...)
	if diags.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	// Refresh the state
	stateRefreshed, found, d := r.read(ctx, state)
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	tflog.Debug(ctx, pp.Sprintf("------State After func read:", stateRefreshed))

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BucketOwnershipControlsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_ownership_controls", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &BucketOwnershipControlsModel{}
		state = &BucketOwnershipControlsModel{}
	)

	// Get current plan and state
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
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
		Implement the resource update here
	*/

	// Set default update timeout
	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	diags.Append(diags...)
	if diags.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	// Update the bucket ownership controls
	d := r.createOrUpdate(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketOwnershipControlsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_ownership_controls", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketOwnershipControlsModel{}

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

	// Set default delete timeout
	deleteTimeout, diags := state.Timeouts.Delete(ctx, defaultDeleteTimeout)
	diags.Append(diags...)
	if diags.HasError() {
		return
	}
	_, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	/*
		Implement the resource deletion here
	*/

	// Delete the bucket ownership controls
	if _, err := r.s3Client.DeleteBucketOwnershipControls(&s3.DeleteBucketOwnershipControlsInput{
		Bucket: state.Bucket.GetPtr(),
	}); err != nil {
		resp.Diagnostics.AddError("Error deleting bucket ownership controls", err.Error())
		return
	}
}

func (r *BucketOwnershipControlsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_ownership_controls", r.client.GetOrgName(), metrics.Import)()

	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketOwnershipControlsResource) read(ctx context.Context, planOrState *BucketOwnershipControlsModel) (stateRefreshed *BucketOwnershipControlsModel, found bool, diags diag.Diagnostics) {
	return genericReadOwnerShipControls(ctx, &genericOwnershipControlsConfig[*BucketOwnershipControlsModel]{
		Client:     r.s3Client,
		BucketName: func() *string { return planOrState.Bucket.GetPtr() },
	}, planOrState)
}

// toAPIOwnershipControlsRules is a custom function to convert the ownership controls rules to API ownership controls rules.
func toAPIOwnershipControlsRules(ctx context.Context, data *BucketOwnershipControlsModel) ([]*s3.OwnershipControlsRule, diag.Diagnostics) {
	var ownershipControlsRules []*s3.OwnershipControlsRule
	rule, diags := data.Rule.Get(ctx)
	if diags.HasError() {
		return nil, diags
	}
	ownershipControlsRules = append(ownershipControlsRules, &s3.OwnershipControlsRule{ObjectOwnership: rule.ObjectOwnership.GetPtr()})
	return ownershipControlsRules, diags
}

// createOrUpdate is a generic custom function to create or update the bucket ownership controls.
func (r *BucketOwnershipControlsResource) createOrUpdate(ctx context.Context, planOrState *BucketOwnershipControlsModel) (diags diag.Diagnostics) {
	// ? Set rules for API
	rules, d := toAPIOwnershipControlsRules(ctx, planOrState)
	if d.HasError() {
		diags.Append(d...)
		return
	}

	tflog.Debug(ctx, pp.Sprintf("------Rules: %v", rules))

	// ? Create the bucket ownership controls
	if _, err := r.s3Client.PutBucketOwnershipControls(&s3.PutBucketOwnershipControlsInput{
		Bucket: planOrState.Bucket.GetPtr(),
		OwnershipControls: &s3.OwnershipControls{
			Rules: rules,
		},
	}); err != nil {
		d.AddError("Error creating or updating bucket ownership controls", err.Error())
		return
	}
	return
}
