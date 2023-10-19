package s3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BucketACLModel struct {
		AccessControlPolicies supertypes.ListNestedObjectValueOf[BucketACLModelAccessControlPolicy] `tfsdk:"access_control_policies"`
		ACL                   supertypes.StringValue                                                `tfsdk:"acl"`
		Bucket                supertypes.StringValue                                                `tfsdk:"bucket"`
		ID                    supertypes.StringValue                                                `tfsdk:"id"`
		Timeouts              timeoutsR.Value                                                       `tfsdk:"timeouts"`
	}

	BucketACLModelDatasource struct {
		AccessControlPolicy supertypes.SingleNestedObjectValueOf[BucketACLModelAccessControlPolicy] `tfsdk:"access_control_policy"`
		ACL                 supertypes.StringValue                                                  `tfsdk:"acl"`
		Bucket              supertypes.StringValue                                                  `tfsdk:"bucket"`
		ID                  supertypes.StringValue                                                  `tfsdk:"id"`
		Timeouts            timeoutsD.Value                                                         `tfsdk:"timeouts"`
	}

	BucketACLModelAccessControlPolicies []*BucketACLModelAccessControlPolicy
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

func (rm *BucketACLModel) Copy() *BucketACLModel {
	x := &BucketACLModel{}
	utils.ModelCopy(rm, x)
	return x
}

// BucketACLModel SetID.
func (rm *BucketACLModel) SetID(id string) {
	rm.ID.Set(id)
}

// BucketACLModel SetACL.
func (rm *BucketACLModel) SetACL(acl string) {
	rm.ACL.Set(acl)
}

// BucketACLModelAccessControlPolicy SetGrants sets the value of the Grants field.
func (rm *BucketACLModelAccessControlPolicy) SetGrants(ctx context.Context, values BucketACLModelGrants) (diags diag.Diagnostics) {
	return rm.Grants.Set(ctx, values)
}

// SetAccessControlPolicy sets the value of the AccessControlPolicy field.
func (rm *BucketACLModel) SetAccessControlPolicies(ctx context.Context, values BucketACLModelAccessControlPolicies) (diags diag.Diagnostics) {
	return rm.AccessControlPolicies.Set(ctx, values)
}

// GetGrant returns the value of the Grants field.
func (rm *BucketACLModelAccessControlPolicy) GetGrants(ctx context.Context) (values BucketACLModelGrants, diags diag.Diagnostics) {
	return rm.Grants.Get(ctx)
}

// GetOwner returns the value of the Owner field.
func (rm *BucketACLModelAccessControlPolicy) GetOwner(ctx context.Context) (values *BucketACLModelOwner, diags diag.Diagnostics) {
	return rm.Owner.Get(ctx)
}

// GetGrant returns the value of the Grant field.
func (rm *BucketACLModelGrants) GetGrant(ctx context.Context, index int) (values *BucketACLModelGrant) {
	return (*rm)[index]
}

// GetGrantee returns the value of the Grantee field.
func (rm *BucketACLModelGrant) GetGrantee(ctx context.Context) (values *BucketACLModelGrantee, diags diag.Diagnostics) {
	return rm.Grantee.Get(ctx)
}

// SetGrantee sets the value of the Grantee field.
func (rm *BucketACLModelGrant) SetGrantee(ctx context.Context, values *BucketACLModelGrantee) (diags diag.Diagnostics) {
	d := rm.Grantee.Set(ctx, values)
	return d
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
		Grants: supertypes.NewSetNestedObjectValueOfNull[BucketACLModelGrant](ctx),    // .NewSetNestedNull(conf.Schema.GetAttributes()["access_control_policy"].(schema.SingleNestedAttribute).GetNestedObject().GetAttributes()["grants"].(schema.SetNestedAttribute).CustomType.(supertypes.SetNestedType).ElemType),
		Owner:  supertypes.NewSingleNestedObjectValueOfNull[BucketACLModelOwner](ctx), // (conf.Schema.GetAttributes()["access_control_policy"].(schema.SingleNestedAttribute).GetNestedObject().GetAttributes()["owner"].(schema.SingleNestedAttribute).CustomType.(supertypes.SingleNestedType).AttrTypes),
	}
}
