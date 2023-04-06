package vm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"

	superschema "github.com/FrangipaneTeam/terraform-plugin-framework-superschema"
)

type VMResourceModelSettingsGuestProperties map[string]string //nolint:revive

func GuestPropertiesSchema() schema.Attribute {
	return schema.MapAttribute{
		MarkdownDescription: "Key/Value settings for guest properties",
		Optional:            true,
		ElementType:         types.StringType,
	}
}

func GuestPropertiesSuperSchema() superschema.Attribute {
	return superschema.MapAttribute{
		Common: &schemaR.MapAttribute{
			MarkdownDescription: "Key/Value settings for guest properties",
			Computed:            true,
			ElementType:         types.StringType,
		},
		Resource: &schemaR.MapAttribute{
			Optional: true,
		},
	}
}

// GuestPropertiesAttrType returns the type map for the guest properties.
func (g *VMResourceModelSettingsGuestProperties) AttrType() attr.Type {
	return types.StringType
}

// ToAttrValue converts a GuestProperties to an attr.Value.
func (g *VMResourceModelSettingsGuestProperties) toAttrValues(_ context.Context) (attrValues map[string]attr.Value) {
	attrValues = make(map[string]attr.Value, len(*g))

	for k, v := range *g {
		attrValues[k] = types.StringValue(v)
	}

	return
}

// ToPlan converts a GuestProperties to a plan.
func (g *VMResourceModelSettingsGuestProperties) ToPlan(ctx context.Context) basetypes.MapValue {
	if g == nil || len(*g) == 0 {
		return types.MapNull(types.StringType)
	}

	return types.MapValueMust(types.StringType, g.toAttrValues(ctx))
}

// GuestPropertiesRead reads the guest properties from a VM.
func (v VM) GuestPropertiesRead() (guestProperties *VMResourceModelSettingsGuestProperties, err error) {
	// get guest properties
	guest, err := v.GetProductSectionList()
	if err != nil {
		return nil, fmt.Errorf("unable to read guest properties: %w", err)
	}

	guestProperties = &VMResourceModelSettingsGuestProperties{}

	for _, guestProperty := range guest.ProductSection.Property {
		if guestProperty.Value != nil {
			(*guestProperties)[guestProperty.Key] = guestProperty.Value.Value
		}
	}

	return
}
