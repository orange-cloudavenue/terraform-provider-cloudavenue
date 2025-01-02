package iam

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type RightModel struct {
	BundleKey     supertypes.StringValue    `tfsdk:"bundle_key"`
	CategoryID    supertypes.StringValue    `tfsdk:"category_id"`
	Description   supertypes.StringValue    `tfsdk:"description"`
	ID            supertypes.StringValue    `tfsdk:"id"`
	ImpliedRights supertypes.SetNestedValue `tfsdk:"implied_rights"`
	Name          supertypes.StringValue    `tfsdk:"name"`
	RightType     supertypes.StringValue    `tfsdk:"right_type"`
}

// * ImpliedRights.
type RightModelImpliedRights []RightModelImpliedRight

// * ImpliedRight.
type RightModelImpliedRight struct {
	ID   supertypes.StringValue `tfsdk:"id"`
	Name supertypes.StringValue `tfsdk:"name"`
}

func NewIAMRight(t any) *RightModel {
	switch x := t.(type) {
	case tfsdk.State:
		return &RightModel{
			BundleKey:     supertypes.NewStringUnknown(),
			CategoryID:    supertypes.NewStringUnknown(),
			Description:   supertypes.NewStringUnknown(),
			ID:            supertypes.NewStringUnknown(),
			ImpliedRights: supertypes.NewSetNestedUnknown(x.Schema.GetAttributes()["implied_rights"].GetType().(supertypes.SetNestedType).ElementType()),
			Name:          supertypes.NewStringNull(),
			RightType:     supertypes.NewStringUnknown(),
		}

	case tfsdk.Plan:
		return &RightModel{
			BundleKey:     supertypes.NewStringUnknown(),
			CategoryID:    supertypes.NewStringUnknown(),
			Description:   supertypes.NewStringUnknown(),
			ID:            supertypes.NewStringUnknown(),
			ImpliedRights: supertypes.NewSetNestedUnknown(x.Schema.GetAttributes()["implied_rights"].GetType().(supertypes.SetNestedType).ElementType()),
			Name:          supertypes.NewStringNull(),
			RightType:     supertypes.NewStringUnknown(),
		}

	case tfsdk.Config:
		return &RightModel{
			BundleKey:     supertypes.NewStringUnknown(),
			CategoryID:    supertypes.NewStringUnknown(),
			Description:   supertypes.NewStringUnknown(),
			ID:            supertypes.NewStringUnknown(),
			ImpliedRights: supertypes.NewSetNestedUnknown(x.Schema.GetAttributes()["implied_rights"].GetType().(supertypes.SetNestedType).ElementType()),
			Name:          supertypes.NewStringNull(),

			RightType: supertypes.NewStringUnknown(),
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *RightModel) Copy() *RightModel {
	x := &RightModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetImpliedRights returns the value of the ImpliedRights field.
func (rm *RightModel) GetImpliedRights(ctx context.Context) (values RightModelImpliedRights, diags diag.Diagnostics) {
	values = make(RightModelImpliedRights, 0)
	d := rm.ImpliedRights.Get(ctx, &values, false)
	return values, d
}
