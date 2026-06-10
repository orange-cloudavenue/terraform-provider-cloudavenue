/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package vdcg

import (
	"context"

	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	sdkv1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/pkg/utils"
)

type networkContextProfileModelSubAttribute struct {
	Type   supertypes.StringValue        `tfsdk:"type"`
	Values supertypes.SetValueOf[string] `tfsdk:"values"`
}

type networkContextProfileModelAttribute struct {
	AppID        supertypes.StringValue                                                     `tfsdk:"app_id"`
	SubAttribute supertypes.ListNestedObjectValueOf[networkContextProfileModelSubAttribute] `tfsdk:"sub_attribute"`
}

type networkContextProfileModel struct {
	ID           supertypes.StringValue                                                  `tfsdk:"id"`
	Name         supertypes.StringValue                                                  `tfsdk:"name"`
	Description  supertypes.StringValue                                                  `tfsdk:"description"`
	VDCGroupID   supertypes.StringValue                                                  `tfsdk:"vdc_group_id"`
	VDCGroupName supertypes.StringValue                                                  `tfsdk:"vdc_group_name"`
	Scope        supertypes.StringValue                                                  `tfsdk:"scope"`
	Attribute    supertypes.ListNestedObjectValueOf[networkContextProfileModelAttribute] `tfsdk:"attribute"`
}

type networkContextProfileModelDatasource struct {
	ID           supertypes.StringValue                                                  `tfsdk:"id"`
	Name         supertypes.StringValue                                                  `tfsdk:"name"`
	Description  supertypes.StringValue                                                  `tfsdk:"description"`
	VDCGroupID   supertypes.StringValue                                                  `tfsdk:"vdc_group_id"`
	VDCGroupName supertypes.StringValue                                                  `tfsdk:"vdc_group_name"`
	Scope        supertypes.StringValue                                                  `tfsdk:"scope"`
	Attribute    supertypes.ListNestedObjectValueOf[networkContextProfileModelAttribute] `tfsdk:"attribute"`
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

func (rm *networkContextProfileModel) fromSDKProfile(ctx context.Context, p *sdkv1.NetworkContextProfile) diag.Diagnostics {
	var diags diag.Diagnostics

	rm.ID.Set(p.ID)
	rm.Name.Set(p.Name)
	rm.Description.Set(p.Description)
	rm.Scope.Set(string(p.Scope))

	diags.Append(rm.Attribute.Set(ctx, attributeBlocksFromSDKProfile(ctx, p, &diags))...)
	return diags
}

// attributeBlocksFromSDKProfile is shared between resource and datasource read paths.
func attributeBlocksFromSDKProfile(ctx context.Context, p *sdkv1.NetworkContextProfile, diags *diag.Diagnostics) []*networkContextProfileModelAttribute {
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
	return attrBlocks
}
