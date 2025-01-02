package s3

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	timeoutsD "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	timeoutsR "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type BucketVersioningConfigurationModel struct {
	Timeouts timeoutsR.Value        `tfsdk:"timeouts"`
	ID       supertypes.StringValue `tfsdk:"id"`
	Bucket   supertypes.StringValue `tfsdk:"bucket"`
	Status   supertypes.StringValue `tfsdk:"status"`
}

type BucketVersioningConfigurationDatasourceModel struct {
	Timeouts timeoutsD.Value        `tfsdk:"timeouts"`
	ID       supertypes.StringValue `tfsdk:"id"`
	Bucket   supertypes.StringValue `tfsdk:"bucket"`
	Status   supertypes.StringValue `tfsdk:"status"`
}

func (rm *BucketVersioningConfigurationModel) Copy() *BucketVersioningConfigurationModel {
	x := &BucketVersioningConfigurationModel{}
	utils.ModelCopy(rm, x)
	return x
}
