package s3

import (
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type BucketModel struct {
	ID         supertypes.StringValue `tfsdk:"id"`
	Name       supertypes.StringValue `tfsdk:"name"`
	ObjectLock supertypes.BoolValue   `tfsdk:"object_lock"`
	Endpoint   supertypes.StringValue `tfsdk:"endpoint"`
}

func (rm *BucketModel) Copy() *BucketModel {
	x := &BucketModel{}
	utils.ModelCopy(rm, x)
	return x
}
