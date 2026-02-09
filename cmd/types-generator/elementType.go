/*
 * SPDX-FileCopyrightText: Copyright (c) 2026 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

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
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type elementType string

const (
	elementTypeBool    elementType = "basetypes.BoolType"
	elementTypeInt64   elementType = "basetypes.Int64Type"
	elementTypeFloat64 elementType = "basetypes.Float64Type"
	elementTypeString  elementType = "basetypes.StringType"
	elementTypeList    elementType = "basetypes.ListType"
	elementTypeSet     elementType = "basetypes.SetType"
	elementTypeMap     elementType = "basetypes.MapType"
)

// New returns a new elementType.
func NewElementType(attribute any) elementType {
	var eType string

	switch x := attribute.(type) {
	case schemaR.SetAttribute:
		eType = x.ElementType.String()
	case schemaR.ListAttribute:
		eType = x.ElementType.String()
	case schemaR.MapAttribute:
		eType = x.ElementType.String()

	case schemaD.SetAttribute:
		eType = x.ElementType.String()
	case schemaD.ListAttribute:
		eType = x.ElementType.String()
	case schemaD.MapAttribute:
		eType = x.ElementType.String()

	case schemaR.StringAttribute, schemaD.StringAttribute:
		return elementTypeString
	case schemaR.BoolAttribute, schemaD.BoolAttribute:
		return elementTypeBool
	case schemaR.Int64Attribute, schemaD.Int64Attribute:
		return elementTypeInt64
	case schemaR.Float64Attribute, schemaD.Float64Attribute:
		return elementTypeFloat64
	}

	switch eType {
	case elementTypeBool.String():
		return elementTypeBool
	case elementTypeInt64.String():
		return elementTypeInt64
	case elementTypeFloat64.String():
		return elementTypeFloat64
	case elementTypeString.String():
		return elementTypeString
	case elementTypeList.String():
		return elementTypeList
	case elementTypeSet.String():
		return elementTypeSet
	case elementTypeMap.String():
		return elementTypeMap
	}

	return ""
}

// String returns the string representation of the elementType.
func (e elementType) String() string {
	return string(e)
}

// ToTerraformType returns the terraform type of the elementType.
func (e elementType) ToTerraformType() string {
	switch e {
	case elementTypeBool:
		return "[]types.Bool"
	case elementTypeInt64:
		return "[]types.Int64"
	case elementTypeFloat64:
		return "[]types.Float64"
	case elementTypeString:
		return "[]types.String"
	}

	return ""
}

// ToGolementType returns the go type of the elementType.
func (e elementType) ToGoType() string {
	switch e {
	case elementTypeBool:
		return "[]bool"
	case elementTypeInt64:
		return "[]int64"
	case elementTypeFloat64:
		return "[]float64"
	case elementTypeString:
		return "[]string"
	}

	return ""
}

// ValidElementType returns true if the elementType is valid.
func ValidElementType(elementType string) bool {
	switch elementType {
	case elementTypeBool.String():
		return true
	case elementTypeInt64.String():
		return true
	case elementTypeFloat64.String():
		return true
	case elementTypeString.String():
		return true
	}

	return false
}
