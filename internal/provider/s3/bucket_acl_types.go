/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BucketACLModel struct {
		AccessControlPolicy supertypes.SingleNestedObjectValueOf[BucketACLModelAccessControlPolicy] `tfsdk:"access_control_policy"`
		ACL                 supertypes.StringValue                                                  `tfsdk:"acl"`
		Bucket              supertypes.StringValue                                                  `tfsdk:"bucket"`
		ID                  supertypes.StringValue                                                  `tfsdk:"id"`
		Timeouts            timeoutsR.Value                                                         `tfsdk:"timeouts"`
	}

	BucketACLModelDatasource struct {
		AccessControlPolicy supertypes.SingleNestedObjectValueOf[BucketACLModelAccessControlPolicy] `tfsdk:"access_control_policy"`
		ACL                 supertypes.StringValue                                                  `tfsdk:"acl"`
		Bucket              supertypes.StringValue                                                  `tfsdk:"bucket"`
		ID                  supertypes.StringValue                                                  `tfsdk:"id"`
		Timeouts            timeoutsD.Value                                                         `tfsdk:"timeouts"`
	}

	// * AccessControlPolicy:
	// Grant field is Optional.
	// Owner field is Require.
	BucketACLModelAccessControlPolicy struct {
		Grants supertypes.SetNestedObjectValueOf[BucketACLModelGrant]    `tfsdk:"grants"`
		Owner  supertypes.SingleNestedObjectValueOf[BucketACLModelOwner] `tfsdk:"owner"`
	}

	// * Grants.
	BucketACLModelGrants []*BucketACLModelGrant

	// * Grant.
	BucketACLModelGrant struct {
		Grantee    supertypes.SingleNestedObjectValueOf[BucketACLModelGrantee] `tfsdk:"grantee"`
		Permission supertypes.StringValue                                      `tfsdk:"permission"`
	}

	// * Owner.
	BucketACLModelOwner struct {
		DisplayName supertypes.StringValue `tfsdk:"display_name"`
		ID          supertypes.StringValue `tfsdk:"id"`
	}

	// * Grantee.
	BucketACLModelGrantee struct {
		DisplayName  supertypes.StringValue `tfsdk:"display_name"`
		EmailAddress supertypes.StringValue `tfsdk:"email_address"`
		ID           supertypes.StringValue `tfsdk:"id"`
		Type         supertypes.StringValue `tfsdk:"type"`
		URI          supertypes.StringValue `tfsdk:"uri"`
	}
)

func (rm *BucketACLModel) Copy() any {
	x := &BucketACLModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *BucketACLModelDatasource) Copy() any {
	x := &BucketACLModelDatasource{}
	utils.ModelCopy(rm, x)
	return x
}

// BucketACLModel SetID.
func (rm *BucketACLModel) SetID(id *string) {
	rm.ID.Set(*id)
}

// BucketACLModelDatasource SetID.
func (rm *BucketACLModelDatasource) SetID(id *string) {
	rm.ID.Set(*id)
}

// SetAccessControlPolicy sets the value of the AccessControlPolicy field.
func (rm *BucketACLModel) SetAccessControlPolicy(ctx context.Context, values *BucketACLModelAccessControlPolicy) (diags diag.Diagnostics) {
	return rm.AccessControlPolicy.Set(ctx, values)
}

// SetAccessControlPolicy sets the value of the AccessControlPolicy field.
func (rm *BucketACLModelDatasource) SetAccessControlPolicy(ctx context.Context, values *BucketACLModelAccessControlPolicy) (diags diag.Diagnostics) {
	return rm.AccessControlPolicy.Set(ctx, values)
}

// NewBucketACLModelGrantee returns a null struct of BucketACLModelGrantee.
func NewBucketACLModelGrantee() *BucketACLModelGrantee {
	return &BucketACLModelGrantee{
		DisplayName:  supertypes.NewStringNull(),
		EmailAddress: supertypes.NewStringNull(),
		ID:           supertypes.NewStringNull(),
		Type:         supertypes.NewStringNull(),
		URI:          supertypes.NewStringNull(),
	}
}

// NewBucketACLModelOwner returns a null struct of BucketACLModelOwner.
func NewBucketACLModelOwner() *BucketACLModelOwner {
	return &BucketACLModelOwner{
		DisplayName: supertypes.NewStringNull(),
		ID:          supertypes.NewStringNull(),
	}
}

// NewBucketACLModelGrant returns a null struct of BucketACLModelGrant.
func NewBucketACLModelGrant(ctx context.Context) *BucketACLModelGrant {
	return &BucketACLModelGrant{
		Grantee:    supertypes.NewSingleNestedObjectValueOfNull[BucketACLModelGrantee](ctx),
		Permission: supertypes.NewStringNull(),
	}
}

// NewBucketACLModelGrants returns a nul struct of BucketACLModelAccessControlPolicy.
func NewBucketACLModelAccessControlPolicy(ctx context.Context) *BucketACLModelAccessControlPolicy {
	return &BucketACLModelAccessControlPolicy{
		Grants: supertypes.NewSetNestedObjectValueOfNull[BucketACLModelGrant](ctx),
		Owner:  supertypes.NewSingleNestedObjectValueOfNull[BucketACLModelOwner](ctx),
	}
}

// Generic struct for resource and datasource.
type readACLResourceDatasource interface {
	*BucketACLModel | *BucketACLModelDatasource
	SetID(*string)
	SetAccessControlPolicy(context.Context, *BucketACLModelAccessControlPolicy) diag.Diagnostics
	Copy() any
}

// Generic type for the datasource.
type readBucketACLGeneric[T readACLResourceDatasource] struct {
	Client     *s3.S3
	BucketName func() *string
}

// Generic function for the datasource.
func genericReadACL[T readACLResourceDatasource](ctx context.Context, config *readBucketACLGeneric[T], planOrState T) (stateRefreshed T, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy().(T)

	bucketACLOutput, err := config.Client.GetBucketAclWithContext(ctx, &s3.GetBucketAclInput{
		Bucket: config.BucketName(),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeNoSuchBucket):
			diags.AddError("Bucket not found", err.Error())
			return nil, false, diags
		default:
			diags.AddError("Error on requesting ACL Bucket", err.Error())
			return nil, true, diags
		}
	}
	if bucketACLOutput == nil {
		diags.AddError("Bucket ACL not found", "bucketACLOutput is nil")
		return nil, true, diags
	}

	// ? Set Owner
	if bucketACLOutput.Owner == nil {
		diags.AddError("Bucket ACL not found", "Owner is nil")
		return nil, true, diags
	}
	owner := NewBucketACLModelOwner()
	owner.DisplayName.SetPtr(bucketACLOutput.Owner.DisplayName)
	owner.ID.SetPtr(bucketACLOutput.Owner.ID)

	// ? Set Grants
	grants := make(BucketACLModelGrants, 0)
	for _, s3Grant := range bucketACLOutput.Grants {
		// ? Set Grant
		grantModel := NewBucketACLModelGrant(ctx)

		// ? Set Permission in grant
		grantModel.Permission.Set(*s3Grant.Permission)

		// ? Set Grantee in grant
		if s3Grant.Grantee == nil {
			diags.AddError("Bucket ACL not found", "Grantee is nil")
			return nil, true, diags
		}
		grantee := NewBucketACLModelGrantee()
		grantee.Type.SetPtr(s3Grant.Grantee.Type)
		grantee.ID.SetPtr(s3Grant.Grantee.ID)
		grantee.DisplayName.SetPtr(s3Grant.Grantee.DisplayName)
		grantee.EmailAddress.SetPtr(s3Grant.Grantee.EmailAddress)
		grantee.URI.SetPtr(s3Grant.Grantee.URI)

		diags.Append(grantModel.Grantee.Set(ctx, grantee)...)
		if diags.HasError() {
			return nil, true, diags
		}

		grants = append(grants, grantModel)
	}

	if len(grants) == 0 {
		diags.AddError("Bucket ACL not found", "Grants or AccessControlPolicy is empty")
		return nil, false, diags
	}

	//  ? Set AccessControlPolicy
	accessControlPolicy := NewBucketACLModelAccessControlPolicy(ctx)
	diags.Append(accessControlPolicy.Grants.Set(ctx, grants)...)
	diags.Append(accessControlPolicy.Owner.Set(ctx, owner)...)
	if diags.HasError() {
		return nil, true, diags
	}

	// Set the accessControlPolicy in the stateRefreshed
	diags.Append(stateRefreshed.SetAccessControlPolicy(ctx, accessControlPolicy)...)
	if diags.HasError() {
		return nil, true, diags
	}

	stateRefreshed.SetID(config.BucketName())

	return stateRefreshed, true, nil
}
