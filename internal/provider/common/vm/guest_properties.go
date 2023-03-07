package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vmware/go-vcloud-director/v2/govcd"
)

type GuestProperties map[string]string

func GuestPropertiesSchema() schema.Attribute {
	return schema.MapAttribute{
		MarkdownDescription: "Key/Value settings for guest properties",
		Optional:            true,
		ElementType:         types.StringType,
	}
}

// GuestPropertiesAttrType returns the type map for the guest properties
func GuestPropertiesAttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"guest_properties": types.MapType{ElemType: types.StringType},
	}
}

// ToAttrValue converts a GuestProperties to an attr.Value
func (g *GuestProperties) ToAttrValue() map[string]attr.Value {
	x := make(map[string]attr.Value)

	if g == nil || len(*g) == 0 {
		return x
	}

	for k, v := range *g {
		x[k] = types.StringValue(v)
	}

	return x

}

// ToPlan converts a GuestProperties to a plan
func (g *GuestProperties) ToPlan() basetypes.MapValue {

	if g == nil || len(*g) == 0 {
		return types.MapNull(types.StringType)
	}

	return types.MapValueMust(types.StringType, g.ToAttrValue())
}

// GuestPropertiesFromPlan converts a terraform plan to a GuestProperties struct.
func GuestPropertiesFromPlan(ctx context.Context, x types.Map) (*GuestProperties, diag.Diagnostics) {
	if x.IsNull() || x.IsUnknown() {
		return &GuestProperties{}, diag.Diagnostics{}
	}

	g := &GuestProperties{}

	d := x.ElementsAs(ctx, g, false)

	return g, d
}

// GuestPropertiesRead reads the guest properties from a VM
func GuestPropertiesRead(vm *govcd.VM) (m *GuestProperties, err error) {

	if m == nil {
		m = &GuestProperties{}
	}

	// get guest properties
	guestProperties, err := vm.GetProductSectionList()
	if err != nil {
		return m, fmt.Errorf("unable to read guest properties: %s", err)
	}

	for _, guestProperty := range guestProperties.ProductSection.Property {
		if guestProperty.Value != nil {
			(*m)[guestProperty.Key] = guestProperty.Value.Value
		}
	}

	return m, nil
}
