package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/metrics"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BucketACLResource{}
	_ resource.ResourceWithConfigure   = &BucketACLResource{}
	_ resource.ResourceWithImportState = &BucketACLResource{}
)

// NewBucketACLResource is a helper function to simplify the provider implementation.
func NewBucketACLResource() resource.Resource {
	return &BucketACLResource{}
}

// BucketACLResource is the resource implementation.
type BucketACLResource struct {
	client   *client.CloudAvenue
	s3Client v1.S3Client
}

// Init Initializes the resource.
func (r *BucketACLResource) Init(ctx context.Context, rm *BucketACLModel) (diags diag.Diagnostics) {
	r.s3Client = r.client.CAVSDK.V1.S3()
	return
}

// Metadata returns the resource type name.
func (r *BucketACLResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + categoryName + "_bucket_acl"
}

// Schema defines the schema for the resource.
func (r *BucketACLResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = bucketACLSchema(ctx).GetResource(ctx)
}

func (r *BucketACLResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *BucketACLResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_acl", r.client.GetOrgName(), metrics.Create)()

	plan := &BucketACLModel{}

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

	// Set default timeouts
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

	// Use generic createOrUpdate function to create the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan, createTimeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	state, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *BucketACLResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer metrics.New("cloudavenue_s3_bucket_acl", r.client.GetOrgName(), metrics.Read)()

	state := &BucketACLModel{}

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

	readTimeout, diags := state.Timeouts.Read(ctx, defaultReadTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
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

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *BucketACLResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_acl", r.client.GetOrgName(), metrics.Update)()

	var (
		plan  = &BucketACLModel{}
		state = &BucketACLModel{}
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

	// Set default timeouts
	updateTimeout, diags := plan.Timeouts.Update(ctx, defaultUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	/*
		Implement the resource update here
	*/

	// Use generic createOrUpdate function to create the resource
	resp.Diagnostics.Append(r.createOrUpdate(ctx, plan, updateTimeout)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use generic read function to refresh the state
	stateRefreshed, _, d := r.read(ctx, plan)
	resp.Diagnostics.Append(d...)
	if d.HasError() {
		return
	}

	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, stateRefreshed)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *BucketACLResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer metrics.New("cloudavenue_s3_bucket_acl", r.client.GetOrgName(), metrics.Delete)()

	state := &BucketACLModel{}

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

	resp.Diagnostics.AddWarning("Note:",
		"\"terraform destroy\" does not delete the S3 Bucket ACL but does remove the resource from Terraform state.",
	)
}

func (r *BucketACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_s3_bucket_acl", r.client.GetOrgName(), metrics.Import)()

	// * Import basic
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketACLResource) read(ctx context.Context, planOrState *BucketACLModel) (stateRefreshed *BucketACLModel, found bool, diags diag.Diagnostics) {
	return genericReadACL(ctx, &readBucketACLGeneric[*BucketACLModel]{
		Client: r.s3Client.S3,
		BucketName: func() *string {
			return planOrState.Bucket.GetPtr()
		},
	}, planOrState)
}

// createOrUpdate is a generic create or update function that can be used by the resource Create and Update functions.
func (r *BucketACLResource) createOrUpdate(ctx context.Context, planOrState *BucketACLModel, timeout time.Duration) (diags diag.Diagnostics) {
	// Set Bucket in input
	input := &s3.PutBucketAclInput{
		Bucket: planOrState.Bucket.GetPtr(),
	}

	if planOrState.ACL.IsKnown() {
		input.ACL = planOrState.ACL.GetPtr()
	} else {
		// ? Get AccessControlPolicy
		accessControlPolicy, d := planOrState.AccessControlPolicy.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		// ? Set Owner
		owner, d := accessControlPolicy.Owner.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		input.AccessControlPolicy = &s3.AccessControlPolicy{
			Owner: &s3.Owner{
				DisplayName: owner.DisplayName.GetPtr(),
				ID:          owner.ID.GetPtr(),
			},
		}

		// ? Set Grants
		if accessControlPolicy.Grants.IsKnown() {
			input.AccessControlPolicy.Grants = make([]*s3.Grant, 0)
			// Get Array of grants
			grants, d := accessControlPolicy.Grants.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}
			// ? Set Grantee and Permission for each grant
			for _, grant := range grants {
				grantee, d := grant.Grantee.Get(ctx)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
				input.AccessControlPolicy.Grants = append(input.AccessControlPolicy.Grants, &s3.Grant{
					Grantee: &s3.Grantee{
						ID:           grantee.ID.GetPtr(),
						DisplayName:  grantee.DisplayName.GetPtr(),
						EmailAddress: grantee.EmailAddress.GetPtr(),
						URI:          grantee.URI.GetPtr(),
						Type:         grantee.Type.GetPtr(),
					},
					Permission: grant.Permission.GetPtr(),
				})
			}
		}
	}

	if _, err := retryWhenAWSErrCodeEquals(ctx, []string{ErrCodeNoSuchBucket, ErrCodeObjectLockConfigurationNotFoundError}, &RetryWhenConfig[*s3.PutBucketAclOutput]{
		Timeout:  timeout,
		Interval: 15 * time.Second,
		Function: func() (*s3.PutBucketAclOutput, error) {
			return r.s3Client.PutBucketAcl(input)
		},
	}); err != nil {
		diags.AddError("Error putting website configuration", err.Error())
		return
	}

	return
}
