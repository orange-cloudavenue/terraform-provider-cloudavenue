package vdc

import (
	"context"
	"fmt"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type GroupModel struct {
	Description supertypes.StringValue `tfsdk:"description"`
	ID          supertypes.StringValue `tfsdk:"id"`
	Name        supertypes.StringValue `tfsdk:"name"`
	Status      supertypes.StringValue `tfsdk:"status"`
	Type        supertypes.StringValue `tfsdk:"type"`
	VDCIds      supertypes.SetValue    `tfsdk:"vdc_ids"`
}

type GroupModelVDCIds []supertypes.StringValue

func NewGroup(t any) *GroupModel {
	switch x := t.(type) {
	case tfsdk.State: //nolint:dupl
		return &GroupModel{
			Description: supertypes.NewStringNull(),
			ID:          supertypes.NewStringUnknown(),
			Name:        supertypes.NewStringNull(),
			Status:      supertypes.NewStringUnknown(),
			Type:        supertypes.NewStringUnknown(),
			VDCIds:      supertypes.NewSetNull(x.Schema.GetAttributes()["vdc_ids"].GetType().(supertypes.SetType).ElementType()),
		}
	case tfsdk.Plan: //nolint:dupl
		return &GroupModel{
			Description: supertypes.NewStringNull(),
			ID:          supertypes.NewStringUnknown(),
			Name:        supertypes.NewStringNull(),
			Status:      supertypes.NewStringUnknown(),
			Type:        supertypes.NewStringUnknown(),
			VDCIds:      supertypes.NewSetNull(x.Schema.GetAttributes()["vdc_ids"].GetType().(supertypes.SetType).ElementType()),
		}
	case tfsdk.Config: //nolint:dupl
		return &GroupModel{
			Description: supertypes.NewStringNull(),
			ID:          supertypes.NewStringUnknown(),
			Name:        supertypes.NewStringNull(),
			Status:      supertypes.NewStringUnknown(),
			Type:        supertypes.NewStringUnknown(),
			VDCIds:      supertypes.NewSetNull(x.Schema.GetAttributes()["vdc_ids"].GetType().(supertypes.SetType).ElementType()),
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

func (rm *GroupModel) Copy() *GroupModel {
	x := &GroupModel{}
	utils.ModelCopy(rm, x)
	return x
}

// GetVDCIds returns the value of the VdcIds field.
func (rm *GroupModel) GetVDCIds(ctx context.Context) (values GroupModelVDCIds, diags diag.Diagnostics) {
	values = make(GroupModelVDCIds, 0)
	d := rm.VDCIds.Get(ctx, &values, false)
	return values, d
}

// Get returns the values.
func (rmVDCIds *GroupModelVDCIds) Get() []string {
	return utils.SuperSliceTypesStringToSliceString(*rmVDCIds)
}

// GetNameOrID returns the name or the id of the resource.
func (rm *GroupModel) GetNameOrID() string {
	if rm.Name.IsKnown() {
		return rm.Name.Get()
	}

	return rm.ID.Get()
}
