/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegw

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

// networkContextProfileModelSubAttribute is a single sub-attribute block.
type networkContextProfileModelSubAttribute struct {
	Type   supertypes.StringValue        `tfsdk:"type"`
	Values supertypes.SetValueOf[string] `tfsdk:"values"`
}

// networkContextProfileModelAttribute is a single attribute block (one App ID + optional sub-attributes).
type networkContextProfileModelAttribute struct {
	AppID        supertypes.StringValue                                                     `tfsdk:"app_id"`
	SubAttribute supertypes.ListNestedObjectValueOf[networkContextProfileModelSubAttribute] `tfsdk:"sub_attribute"`
}

// networkContextProfileModel is the resource model.
type networkContextProfileModel struct {
	ID              supertypes.StringValue                                                  `tfsdk:"id"`
	Name            supertypes.StringValue                                                  `tfsdk:"name"`
	Description     supertypes.StringValue                                                  `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue                                                  `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                                  `tfsdk:"edge_gateway_name"`
	Scope           supertypes.StringValue                                                  `tfsdk:"scope"`
	Attribute       supertypes.ListNestedObjectValueOf[networkContextProfileModelAttribute] `tfsdk:"attribute"`
}

// networkContextProfileModelDatasource is the datasource model.
type networkContextProfileModelDatasource struct {
	ID              supertypes.StringValue                                                  `tfsdk:"id"`
	Name            supertypes.StringValue                                                  `tfsdk:"name"`
	Description     supertypes.StringValue                                                  `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue                                                  `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                                  `tfsdk:"edge_gateway_name"`
	Scope           supertypes.StringValue                                                  `tfsdk:"scope"`
	Attribute       supertypes.ListNestedObjectValueOf[networkContextProfileModelAttribute] `tfsdk:"attribute"`
}

func (rm *networkContextProfileModel) Copy() *networkContextProfileModel {
	x := &networkContextProfileModel{}
	utils.ModelCopy(rm, x)
	return x
}

func (rm *networkContextProfileModelDatasource) Copy() *networkContextProfileModelDatasource {
	x := &networkContextProfileModelDatasource{}
	utils.ModelCopy(rm, x)
	return x
}

// toSDKProfile converts the resource model to the SDK NetworkContextProfile.
func (rm *networkContextProfileModel) toSDKProfile(ctx context.Context) (*sdkv1.NetworkContextProfile, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrBlocks, d := rm.Attribute.Get(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}

	attrs := make([]sdkv1.NetworkContextProfileAttribute, 0, len(attrBlocks))
	for _, block := range attrBlocks {
		subAttrBlocks, d := block.SubAttribute.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		subAttrs := make([]sdkv1.NetworkContextProfileSubAttribute, 0, len(subAttrBlocks))
		for _, sub := range subAttrBlocks {
			values, d := sub.Values.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
			subAttrs = append(subAttrs, sdkv1.NetworkContextProfileSubAttribute{
				Type:   sdkv1.NetworkContextProfileSubAttributeType(sub.Type.Get()),
				Values: values,
			})
		}

		attrs = append(attrs, sdkv1.NetworkContextProfileAttribute{
			Type:          sdkv1.NetworkContextProfileAttributeTypeAppID,
			Values:        []string{block.AppID.Get()},
			SubAttributes: subAttrs,
		})
	}

	return &sdkv1.NetworkContextProfile{
		ID:          rm.ID.Get(),
		Name:        rm.Name.Get(),
		Description: rm.Description.Get(),
		Attributes:  attrs,
	}, diags
}

// fromSDKProfile populates the resource model from the SDK NetworkContextProfile.
func (rm *networkContextProfileModel) fromSDKProfile(ctx context.Context, p *sdkv1.NetworkContextProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	rm.ID.Set(p.ID)
	rm.Name.Set(p.Name)
	rm.Description.Set(p.Description)
	rm.Scope.Set(string(p.Scope))

	attrBlocks, d := attributesFromSDKProfile(ctx, p)
	diags.Append(d...)
	diags.Append(rm.Attribute.Set(ctx, attrBlocks)...)
	return diags
}

// attributesFromSDKProfile is a shared helper for the datasource read path.
func attributesFromSDKProfile(ctx context.Context, p *sdkv1.NetworkContextProfile) ([]*networkContextProfileModelAttribute, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrBlocks := make([]*networkContextProfileModelAttribute, 0)
	for _, attr := range p.Attributes {
		if attr.Type != sdkv1.NetworkContextProfileAttributeTypeAppID {
			continue
		}
		for _, appID := range attr.Values {
			subAttrBlocks := make([]*networkContextProfileModelSubAttribute, 0, len(attr.SubAttributes))
			for _, s := range attr.SubAttributes {
				sub := &networkContextProfileModelSubAttribute{
					Type:   supertypes.NewStringNull(),
					Values: supertypes.NewSetValueOfNull[string](ctx),
				}
				sub.Type.Set(string(s.Type))
				diags.Append(sub.Values.Set(ctx, s.Values)...)
				subAttrBlocks = append(subAttrBlocks, sub)
			}

			block := &networkContextProfileModelAttribute{
				AppID:        supertypes.NewStringNull(),
				SubAttribute: supertypes.NewListNestedObjectValueOfNull[networkContextProfileModelSubAttribute](ctx),
			}
			block.AppID.Set(appID)
			diags.Append(block.SubAttribute.Set(ctx, subAttrBlocks)...)
			attrBlocks = append(attrBlocks, block)
		}
	}

	return attrBlocks, diags
}
