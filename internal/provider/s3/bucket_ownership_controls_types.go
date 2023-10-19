package s3

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/k0kubun/pp"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type (
	BucketOwnershipControlsModel struct {
		Bucket   supertypes.StringValue                     `tfsdk:"bucket"`
		ID       supertypes.StringValue                     `tfsdk:"id"`
		Rule     supertypes.SingleNestedObjectValueOf[rule] `tfsdk:"rule"`
		Timeouts timeoutsR.Value                            `tfsdk:"timeouts"`
	}

	BucketOwnershipControlsDataSourceModel struct {
		Bucket   supertypes.StringValue                     `tfsdk:"bucket"`
		ID       supertypes.StringValue                     `tfsdk:"id"`
		Rule     supertypes.SingleNestedObjectValueOf[rule] `tfsdk:"rule"`
		Timeouts timeoutsD.Value                            `tfsdk:"timeouts"`
	}

	rule struct {
		ObjectOwnership supertypes.StringValue `tfsdk:"object_ownership"`
	}
)

func (rm *BucketOwnershipControlsModel) Copy() any {
	x := &BucketOwnershipControlsModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *BucketOwnershipControlsDataSourceModel) Copy() any {
	x := &BucketOwnershipControlsDataSourceModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *BucketOwnershipControlsModel) SetRule(ctx context.Context, rule *rule) diag.Diagnostics {
	return rm.Rule.Set(ctx, rule)
}

func (rm *BucketOwnershipControlsDataSourceModel) SetRule(ctx context.Context, rule *rule) diag.Diagnostics {
	return rm.Rule.Set(ctx, rule)
}

func (rm *BucketOwnershipControlsModel) SetID(id *string) {
	rm.ID.SetPtr(id)
}

func (rm *BucketOwnershipControlsDataSourceModel) SetID(id *string) {
	rm.ID.SetPtr(id)
}

// genericOwnershipControls is a generic interface to read the ownership controls.
type genericOwnershipControls interface {
	*BucketOwnershipControlsModel | *BucketOwnershipControlsDataSourceModel
	Copy() any
	SetRule(ctx context.Context, rule *rule) diag.Diagnostics
	SetID(id *string)
}

// genericOwnershipControlsConfig is a generic struct.
type genericOwnershipControlsConfig[T genericOwnershipControls] struct {
	Client     *s3.S3
	BucketName func() *string
}

// genericReadOwnerShipControls is a generic function to read the ownership controls.
func genericReadOwnerShipControls[T genericOwnershipControls](ctx context.Context, config *genericOwnershipControlsConfig[T], planOrState T) (stateRefreshed T, found bool, diags diag.Diagnostics) {
	stateRefreshed = planOrState.Copy().(T)

	// Get the API bucket ownership controls
	APIOwnershipControl, err := config.Client.GetBucketOwnershipControls(&s3.GetBucketOwnershipControlsInput{
		Bucket: config.BucketName(),
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), ErrCodeOwnershipControlsNotFoundError):
			diags.AddError("Error retrieving bucket ownership controls", err.Error())
			return nil, false, diags
		case strings.Contains(err.Error(), s3.ErrCodeNoSuchBucket):
			diags.AddError("Error bucket not found", err.Error())
			return nil, false, diags
		default:
			diags.AddError("Error retrieving bucket", err.Error())
			return nil, false, diags
		}
	}

	tflog.Debug(ctx, pp.Sprintf("APIOwnershipControl: %v", APIOwnershipControl))

	// ? Set rule from API
	rule := &rule{}
	if len(APIOwnershipControl.OwnershipControls.Rules) > 0 {
		rule.ObjectOwnership.Set(APIOwnershipControl.OwnershipControls.String())
	}
	if diags = stateRefreshed.SetRule(ctx, rule); diags.HasError() {
		return nil, true, diags
	}

	stateRefreshed.SetID(config.BucketName())

	return stateRefreshed, true, nil
}
