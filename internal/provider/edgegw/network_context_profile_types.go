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

// networkContextProfileModelSubAttribute is a single sub-attribute block within app_id.
type networkContextProfileModelSubAttribute struct {
	Type   supertypes.StringValue        `tfsdk:"type"`
	Values supertypes.SetValueOf[string] `tfsdk:"values"`
}

// networkContextProfileModelAppID holds the APP_ID attribute block.
type networkContextProfileModelAppID struct {
	Values       supertypes.SetValueOf[string]                                              `tfsdk:"values"`
	SubAttribute supertypes.ListNestedObjectValueOf[networkContextProfileModelSubAttribute] `tfsdk:"sub_attributes"`
}

// networkContextProfileModelDomainName holds the DOMAIN_NAME attribute block.
type networkContextProfileModelDomainName struct {
	Values supertypes.SetValueOf[string] `tfsdk:"values"`
}

// networkContextProfileModel is the resource model.
type networkContextProfileModel struct {
	ID              supertypes.StringValue                                                `tfsdk:"id"`
	Name            supertypes.StringValue                                                `tfsdk:"name"`
	Description     supertypes.StringValue                                                `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue                                                `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                                `tfsdk:"edge_gateway_name"`
	Scope           supertypes.StringValue                                                `tfsdk:"scope"`
	AppID           supertypes.SingleNestedObjectValueOf[networkContextProfileModelAppID] `tfsdk:"app_id"`
}

// networkContextProfileModelDatasource is the datasource model.
type networkContextProfileModelDatasource struct {
	ID              supertypes.StringValue                                                     `tfsdk:"id"`
	Name            supertypes.StringValue                                                     `tfsdk:"name"`
	Description     supertypes.StringValue                                                     `tfsdk:"description"`
	EdgeGatewayID   supertypes.StringValue                                                     `tfsdk:"edge_gateway_id"`
	EdgeGatewayName supertypes.StringValue                                                     `tfsdk:"edge_gateway_name"`
	Scope           supertypes.StringValue                                                     `tfsdk:"scope"`
	AppID           supertypes.SingleNestedObjectValueOf[networkContextProfileModelAppID]      `tfsdk:"app_id"`
	DomainName      supertypes.SingleNestedObjectValueOf[networkContextProfileModelDomainName] `tfsdk:"domain_name"`
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

	attrs := make([]sdkv1.NetworkContextProfileAttribute, 0, 2)

	// APP_ID block
	if rm.AppID.IsKnown() {
		appIDBlock := rm.AppID.DiagsGet(ctx, diags)
		if diags.HasError() {
			return nil, diags
		}

		values, d := appIDBlock.Values.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		subAttrBlocks, d := appIDBlock.SubAttribute.Get(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}

		subAttrs := make([]sdkv1.NetworkContextProfileSubAttribute, 0, len(subAttrBlocks))
		for _, sub := range subAttrBlocks {
			subValues, d := sub.Values.Get(ctx)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
			subAttrs = append(subAttrs, sdkv1.NetworkContextProfileSubAttribute{
				Type:   sdkv1.NetworkContextProfileSubAttributeType(sub.Type.Get()),
				Values: subValues,
			})
		}

		attrs = append(attrs, sdkv1.NetworkContextProfileAttribute{
			Type:          sdkv1.NetworkContextProfileAttributeTypeAppID,
			Values:        values,
			SubAttributes: subAttrs,
		})
	}

	// TODO: DOMAIN_NAME attribute support is not yet implemented.
	// Adding or editing DOMAIN_NAME values requires elevated permissions that are not available through
	// the VCD API directly. The Cerberus middleware (Cloud Avenue platform API) is expected to act as
	// the intermediary for these operations, but this has not been implemented on their side yet.
	// Once Cerberus exposes DOMAIN_NAME management, this block should be wired up.
	// The read path (attributesFromSDKProfile) already handles DOMAIN_NAME for existing SYSTEM profiles.

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

	appIDBlock, _, d := attributesFromSDKProfile(ctx, p)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	diags.Append(rm.AppID.Set(ctx, appIDBlock)...)
	return diags
}

// fromSDKProfile populates the datasource model from the SDK NetworkContextProfile.
func (dm *networkContextProfileModelDatasource) fromSDKProfile(ctx context.Context, p *sdkv1.NetworkContextProfile, edgeGWID, edgeGWName string) diag.Diagnostics {
	var diags diag.Diagnostics

	dm.ID.Set(p.ID)
	dm.Name.Set(p.Name)
	dm.Description.Set(p.Description)
	dm.Scope.Set(string(p.Scope))
	dm.EdgeGatewayID.Set(edgeGWID)
	dm.EdgeGatewayName.Set(edgeGWName)

	appIDBlock, domainBlock, d := attributesFromSDKProfile(ctx, p)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	diags.Append(dm.AppID.Set(ctx, appIDBlock)...)
	diags.Append(dm.DomainName.Set(ctx, domainBlock)...)
	return diags
}

// attributesFromSDKProfile converts SDK profile attributes into the two optional model blocks.
// Returns nil for a block if that attribute type is not present in the profile.
func attributesFromSDKProfile(ctx context.Context, p *sdkv1.NetworkContextProfile) (*networkContextProfileModelAppID, *networkContextProfileModelDomainName, diag.Diagnostics) {
	var diags diag.Diagnostics

	var appIDResult *networkContextProfileModelAppID
	var domainResult *networkContextProfileModelDomainName

	for _, attr := range p.Attributes {
		switch attr.Type {
		case sdkv1.NetworkContextProfileAttributeTypeAppID:
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

			block := &networkContextProfileModelAppID{
				Values:       supertypes.NewSetValueOfNull[string](ctx),
				SubAttribute: supertypes.NewListNestedObjectValueOfNull[networkContextProfileModelSubAttribute](ctx),
			}
			diags.Append(block.Values.Set(ctx, attr.Values)...)
			if len(subAttrBlocks) > 0 {
				diags.Append(block.SubAttribute.Set(ctx, subAttrBlocks)...)
			}
			appIDResult = block

		case sdkv1.NetworkContextProfileAttributeTypeDomainName:
			block := &networkContextProfileModelDomainName{
				Values: supertypes.NewSetValueOfNull[string](ctx),
			}
			diags.Append(block.Values.Set(ctx, attr.Values)...)
			domainResult = block

			// Unknown attribute types silently skipped for forward compatibility.
		}
	}

	return appIDResult, domainResult, diags
}
