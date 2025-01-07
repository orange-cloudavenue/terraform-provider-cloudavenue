/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package main

import (
	"strings"

	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func baseTypesToType(baseType string) string {
	return strings.Replace(baseType, "basetypes.", "types.", 1)
}

// New returns a new elementType.
func NewAttributeType(attribute any) string {
	switch x := attribute.(type) {
	// Schema Resource
	case schemaR.SetAttribute:
		return "types.SetType{ElemType:" + baseTypesToType(x.ElementType.String()) + "}"
	case schemaR.ListAttribute:
		return "types.ListType{ElemType:" + baseTypesToType(x.ElementType.String()) + "}"
	case schemaR.MapAttribute:
		return "types.MapType{ElemType:" + baseTypesToType(x.ElementType.String()) + "}"

	// Schema DataSource
	case schemaD.SetAttribute:
		return "types.SetType{ElemType:" + baseTypesToType(x.ElementType.String()) + "}"
	case schemaD.ListAttribute:
		return "types.ListType{ElemType:" + baseTypesToType(x.ElementType.String()) + "}"
	case schemaD.MapAttribute:
		return "types.ListType{ElemType:" + baseTypesToType(x.ElementType.String()) + "}"

	case schemaR.StringAttribute, schemaD.StringAttribute:
		return "types.StringType"
	case schemaR.BoolAttribute, schemaD.BoolAttribute:
		return "types.BoolType"
	case schemaR.Int64Attribute, schemaD.Int64Attribute:
		return "types.Int64Type"
	case schemaR.Float64Attribute, schemaD.Float64Attribute:
		return "types.Float64Type"
	}

	return ""
}
