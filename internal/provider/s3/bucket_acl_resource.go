package s3

import (
	"context"
	"fmt"
	"strings"

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
	s3Client *s3.S3
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

	tflog.Debug(ctx, pp.Sprint("=====> Before reading Plan: %s", plan))

	// Retrieve values from plan
	resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, pp.Sprint("=====> input.plan: %s", plan))

	// Init the resource
	resp.Diagnostics.Append(r.Init(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default timeouts
	createTimeout, diags := plan.Timeouts.Read(ctx, defaultCreateTimeout)
	diags.Append(diags...)
	if diags.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	/*
		Implement the resource creation logic here.
	*/
	var err error

	// Set Bucket in input
	input := &s3.PutBucketAclInput{
		Bucket: plan.Bucket.GetPtr(),
	}

	if plan.ACL.IsKnown() {
		input.ACL = plan.ACL.GetPtr()
		tflog.Debug(ctx, pp.Sprint("=====> input.ACL: %s", input.ACL))
	} else {
		// ? Get AccessControlPolicy
		accessControlPolicies, diags := plan.AccessControlPolicies.Get(ctx)
		accessControlPolicy := accessControlPolicies[0]
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// ? Set Owner
		owner, diags := accessControlPolicy.GetOwner(ctx)
		// owner := NewBucketACLModelOwner()
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
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
			grants, diags := accessControlPolicy.GetGrants(ctx)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			// ? Set Grantee / Permission for each grant
			for i := 0; i < len(grants); i++ {
				grant := grants.GetGrant(ctx, i)
				grantee, diags := grant.GetGrantee(ctx)
				if diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
				Grantee := &s3.Grantee{}
				if grantee.ID.IsKnown() {
					Grantee.ID = grantee.ID.GetPtr()
				}
				if grantee.DisplayName.IsKnown() {
					Grantee.DisplayName = grantee.DisplayName.GetPtr()
				}
				if grantee.EmailAddress.IsKnown() {
					Grantee.EmailAddress = grantee.EmailAddress.GetPtr()
				}
				if grantee.URI.IsKnown() {
					Grantee.URI = grantee.URI.GetPtr()
				}
				Grantee.Type = grantee.Type.GetPtr()
				input.AccessControlPolicy.Grants = append(input.AccessControlPolicy.Grants, &s3.Grant{
					Grantee:    Grantee,
					Permission: grants[i].Permission.GetPtr(),
				})
			}
		}
	}

	_, err = r.s3Client.PutBucketAclWithContext(ctx, input)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchBucket):
			resp.Diagnostics.AddError("Bucket not found", err.Error())
			return
		case strings.Contains(err.Error(), ErrCodeObjectLockConfigurationNotFoundError):
			resp.Diagnostics.AddError("Bucket object lock configuration was not found", err.Error())
			return
		default:
			resp.Diagnostics.AddError("Bucket ACL not found", err.Error())
			return
		}
	}

	// Use generic read function to refresh the state
	state, _, d := r.read(ctx, plan)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
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

	/*
		Implement the resource update here
	*/
	input := &s3.PutBucketAclInput{
		Bucket: plan.Bucket.GetPtr(),
	}

	if !plan.ACL.Equal(state.ACL) {
		input.ACL = plan.ACL.GetPtr()
	}

	if !plan.AccessControlPolicies.Equal(state.AccessControlPolicies) {
		// ? Get AccessControlPolicy
		accessControlPolicies, diags := plan.AccessControlPolicies.Get(ctx)
		accessControlPolicy := accessControlPolicies[0]
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		// ? Set Owner
		owner, diags := accessControlPolicy.GetOwner(ctx)
		// owner := NewBucketACLModelOwner()
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
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
			grants, diags := accessControlPolicy.GetGrants(ctx)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			// ? Set Grantee / Permission for each grant
			for i := 0; i < len(grants); i++ {
				grant := grants.GetGrant(ctx, i)
				grantee, diags := grant.GetGrantee(ctx)
				if diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
				Grantee := &s3.Grantee{}
				if grantee.ID.IsKnown() {
					Grantee.ID = grantee.ID.GetPtr()
				}
				if grantee.DisplayName.IsKnown() {
					Grantee.DisplayName = grantee.DisplayName.GetPtr()
				}
				if grantee.EmailAddress.IsKnown() {
					Grantee.EmailAddress = grantee.EmailAddress.GetPtr()
				}
				if grantee.URI.IsKnown() {
					Grantee.URI = grantee.URI.GetPtr()
				}
				Grantee.Type = grantee.Type.GetPtr()
				input.AccessControlPolicy.Grants = append(input.AccessControlPolicy.Grants, &s3.Grant{
					Grantee:    Grantee,
					Permission: grants[i].Permission.GetPtr(),
				})
			}
		}
	}

	_, err := r.s3Client.PutBucketAclWithContext(ctx, input)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchBucket):
			resp.Diagnostics.AddError("Bucket not found", err.Error())
			return
		case strings.Contains(err.Error(), ErrCodeObjectLockConfigurationNotFoundError):
			resp.Diagnostics.AddError("Bucket object lock configuration was not found", err.Error())
			return
		default:
			resp.Diagnostics.AddError("Bucket ACL not found", err.Error())
			return
		}
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
	resp.State.RemoveResource(ctx)
}

func (r *BucketACLResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer metrics.New("cloudavenue_backup", r.client.GetOrgName(), metrics.Import)()

	// * Import basic
	resource.ImportStatePassthroughID(ctx, path.Root("bucket"), req, resp)

	// * Import with custom logic
	// idParts := strings.Split(req.ID, ".")

	// if len(idParts) != 2 {
	// 	resp.Diagnostics.AddError(
	// 		"Unexpected Import Identifier",
	// 		fmt.Sprintf("Expected import identifier with format: xx.xx. Got: %q", req.ID),
	// 	)
	// 	return
	// }

	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("xx"), var1)...)
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("xx"), var2)...)
}

// * CustomFuncs

// read is a generic read function that can be used by the resource Create, Read and Update functions.
func (r *BucketACLResource) read(ctx context.Context, planOrState *BucketACLModel) (stateRefreshed *BucketACLModel, found bool, diags diag.Diagnostics) {
	// stateRefreshed is commented because the Copy function is not before run the types generator
	stateRefreshed = planOrState.Copy()

	/*
		Implement the resource read here
	*/

	bucketACLOutput, err := r.s3Client.GetBucketAclWithContext(ctx, &s3.GetBucketAclInput{
		Bucket: planOrState.Bucket.GetPtr(),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchBucket):
			diags.AddError("Bucket not found", err.Error())
			return nil, false, diags
		default:
			diags.AddError("Bucket ACL not found", err.Error())
			return nil, false, diags
		}
	}
	if bucketACLOutput == nil {
		diags.AddError("Bucket ACL not found", "bucketACLOutput is nil")
		return nil, false, diags
	}

	// If no grants in the bucket ACL returned, return the stateRefreshed
	// if len(bucketACLOutput.Grants) == 0 {
	// 	return stateRefreshed, true, nil
	// }

	// ? Set Owner
	s3Owner := bucketACLOutput.Owner
	if s3Owner == nil {
		diags.AddError("Bucket ACL not found", "Owner is nil")
		return nil, false, diags
	}
	owner := NewBucketACLModelOwner()
	owner.DisplayName.Set(*s3Owner.DisplayName)
	owner.ID.Set(*s3Owner.ID)

	// ? Set Grants
	grants := make(BucketACLModelGrants, 0)
	for _, s3Grant := range bucketACLOutput.Grants {
		// ? Set Grant
		grantModel := NewBucketACLModelGrant(ctx)

		// ? Set Permission in grant
		grantModel.Permission.Set(*s3Grant.Permission)

		// ? Set Grantee in grant
		s3Grantee := s3Grant.Grantee
		if s3Grantee == nil {
			diags.AddError("Bucket ACL not found", "Grantee is nil")
			return nil, false, diags
		}
		grantee := NewBucketACLModelGrantee()
		grantee.Type.SetPtr(s3Grantee.Type)
		if s3Grantee.ID != nil {
			grantee.ID.Set(*s3Grantee.ID)
		}
		if s3Grantee.DisplayName != nil {
			grantee.DisplayName.Set(*s3Grantee.DisplayName)
		}
		if s3Grantee.EmailAddress != nil {
			grantee.EmailAddress.Set(*s3Grantee.EmailAddress)
		}
		if s3Grantee.URI != nil {
			grantee.URI.SetPtr(s3Grantee.URI)
		}

		if d := grantModel.SetGrantee(ctx, grantee); d.HasError() {
			diags.Append(d...)
			return nil, false, diags
		}
		grants = append(grants, grantModel)
	}
	tflog.Debug(ctx, pp.Sprint("=====> grants: %s", grants))

	//  ? Set AccessControlPolicy
	accessControlPolicy := NewBucketACLModelAccessControlPolicy(ctx)
	accessControlPolicy.Grants.Set(ctx, grants)
	accessControlPolicy.Owner.Set(ctx, owner)

	tflog.Debug(ctx, pp.Sprint("accessControlPolicy: %s", accessControlPolicy))

	// Set the accessControlPolicy in the stateRefreshed
	slice := make(BucketACLModelAccessControlPolicies, 0)
	if d := stateRefreshed.AccessControlPolicies.Set(ctx, append(slice, accessControlPolicy)); d.HasError() {
		diags.Append(d...)
		return nil, false, diags
	}

	if !stateRefreshed.ID.IsKnown() {
		stateRefreshed.ID.Set(planOrState.Bucket.Get())
	}

	return stateRefreshed, true, nil
}
