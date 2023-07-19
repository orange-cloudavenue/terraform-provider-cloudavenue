package main

type schemaType string

const (
	schemaTypeBool         schemaType = "schema.BoolAttribute"
	schemaTypeInt64        schemaType = "schema.Int64Attribute"
	schemaTypeFloat64      schemaType = "schema.Float64Attribute"
	schemaTypeString       schemaType = "schema.StringAttribute"
	schemaTypeList         schemaType = "schema.ListAttribute"
	schemaTypeSet          schemaType = "schema.SetAttribute"
	schemaTypeMap          schemaType = "schema.MapAttribute"
	schemaTypeListNested   schemaType = "schema.ListNestedAttribute"
	schemaTypeSetNested    schemaType = "schema.SetNestedAttribute"
	schemaTypeMapNested    schemaType = "schema.MapNestedAttribute"
	schemaTypeSingleNested schemaType = "schema.SingleNestedAttribute"
)

// New returns a new schemaType.
func NewSchemaType(schemaType string) schemaType {
	switch schemaType {
	case schemaTypeBool.String():
		return schemaTypeBool
	case schemaTypeInt64.String():
		return schemaTypeInt64
	case schemaTypeFloat64.String():
		return schemaTypeFloat64
	case schemaTypeString.String():
		return schemaTypeString
	case schemaTypeList.String():
		return schemaTypeList
	case schemaTypeSet.String():
		return schemaTypeSet
	case schemaTypeMap.String():
		return schemaTypeMap
	case schemaTypeListNested.String():
		return schemaTypeListNested
	case schemaTypeSetNested.String():
		return schemaTypeSetNested
	case schemaTypeMapNested.String():
		return schemaTypeMapNested
	case schemaTypeSingleNested.String():
		return schemaTypeSingleNested
	}

	return ""
}

// String returns the string representation of the schemaType.
func (s schemaType) String() string {
	return string(s)
}

// ToTerraformType returns the terraform type of the schemaType.
func (s schemaType) ToTerraformValue() string {
	switch s {
	case schemaTypeBool:
		return "supertypes.BoolValue"
	case schemaTypeInt64:
		return "supertypes.Int64Value"
	case schemaTypeFloat64:
		return "supertypes.Float64Value"
	case schemaTypeString:
		return "supertypes.StringValue"
	case schemaTypeList:
		return "supertypes.ListValue"
	case schemaTypeListNested:
		return "supertypes.ListNestedValue"
	case schemaTypeSet:
		return "supertypes.SetValue"
	case schemaTypeSetNested:
		return "supertypes.SetNestedValue"
	case schemaTypeMap:
		return "supertypes.MapValue"
	case schemaTypeMapNested:
		return "supertypes.MapNestedValue"
	case schemaTypeSingleNested:
		return "supertypes.SingleNestedValue"
	}

	return ""
}

// ToTerraformType returns the terraform type of the schemaType.
func (s schemaType) ToTerraformType() string {
	switch s {
	case schemaTypeBool:
		return "supertypes.BoolType"
	case schemaTypeInt64:
		return "supertypes.Int64Type"
	case schemaTypeFloat64:
		return "supertypes.Float64Type"
	case schemaTypeString:
		return "supertypes.StringType"
	case schemaTypeList:
		return "supertypes.ListType"
	case schemaTypeListNested:
		return "supertypes.ListNestedType"
	case schemaTypeSet:
		return "supertypes.SetType"
	case schemaTypeSetNested:
		return "supertypes.SetNestedType"
	case schemaTypeMap:
		return "supertypes.MapType"
	case schemaTypeMapNested:
		return "supertypes.MapNestedType"
	case schemaTypeSingleNested:
		return "supertypes.SingleNestedType"
	}

	return ""
}

// ToBaseTypeValue returns the base type of the elementType.
func (s schemaType) ToBaseTypeValue() string {
	switch s {
	case schemaTypeBool:
		return "basetypes.BoolValue"
	case schemaTypeInt64:
		return "basetypes.Int64Value"
	case schemaTypeFloat64:
		return "basetypes.Float64Value"
	case schemaTypeString:
		return "basetypes.StringValue"
	case schemaTypeList, schemaTypeListNested:
		return "basetypes.ListValue"
	case schemaTypeSet, schemaTypeSetNested:
		return "basetypes.SetValue"
	case schemaTypeMap, schemaTypeMapNested:
		return "basetypes.MapValue"
	}

	return ""
}

// ToFuncNull returns the funcNull of the elementType.
func (s schemaType) ToFuncNull() string {
	switch s {
	case schemaTypeBool:
		return "supertypes.NewBoolNull"
	case schemaTypeInt64:
		return "supertypes.NewInt64Null"
	case schemaTypeFloat64:
		return "supertypes.NewFloat64Null"
	case schemaTypeString:
		return "supertypes.NewStringNull"
	case schemaTypeList:
		return "supertypes.NewListNull"
	case schemaTypeListNested:
		return "supertypes.NewListNestedNull"
	case schemaTypeSet:
		return "supertypes.NewSetNull"
	case schemaTypeSetNested:
		return "supertypes.NewSetNestedNull"
	case schemaTypeMap:
		return "supertypes.NewMapNull"
	case schemaTypeMapNested:
		return "supertypes.NewMapNestedNull"
	case schemaTypeSingleNested:
		return "supertypes.NewSingleNestedNull"
	}

	return ""
}

// ToFuncUnkown returns the funcUnkown of the elementType.
func (s schemaType) ToFuncUnkown() string {
	switch s {
	case schemaTypeBool:
		return "supertypes.NewBoolUnknown"
	case schemaTypeInt64:
		return "supertypes.NewInt64Unknown"
	case schemaTypeFloat64:
		return "supertypes.NewFloat64Unknown"
	case schemaTypeString:
		return "supertypes.NewStringUnknown"
	case schemaTypeList:
		return "supertypes.NewListUnknown"
	case schemaTypeListNested:
		return "supertypes.NewListNestedUnknown"
	case schemaTypeSet:
		return "supertypes.NewSetUnknown"
	case schemaTypeSetNested:
		return "supertypes.NewSetNestedUnknown"
	case schemaTypeMap:
		return "supertypes.NewMapUnknown"
	case schemaTypeMapNested:
		return "supertypes.NewMapNestedUnknown"
	case schemaTypeSingleNested:
		return "supertypes.NewSingleNestedUnknown"
	}

	return ""
}

// ToValueFrom returns the valueFrom of the elementType.
func (s schemaType) ToValueFrom() string {
	switch s {
	case schemaTypeList, schemaTypeListNested:
		return "types.ListValueFrom"
	case schemaTypeSet, schemaTypeSetNested:
		return "types.SetValueFrom"
	case schemaTypeMap, schemaTypeMapNested:
		return "types.MapValueFrom"
	}

	return ""
}

// Valid returns true if the schemaType is valid.
func ValidSchemaType(schemaType string) bool {
	switch schemaType {
	case schemaTypeBool.String(), schemaTypeInt64.String(), schemaTypeFloat64.String(), schemaTypeString.String(), schemaTypeList.String(), schemaTypeSet.String(), schemaTypeMap.String(), schemaTypeListNested.String(), schemaTypeSetNested.String(), schemaTypeMapNested.String():
		return true
	}

	return false
}

// IsArray returns true if the schemaType is an array.
func IsArray(schemaType string) bool {
	return schemaType == schemaTypeList.String() || schemaType == schemaTypeSet.String() || schemaType == schemaTypeMap.String()
}

// IsNested returns true if the schemaType is nested.
func IsNested(schemaType string) bool {
	return schemaType == schemaTypeListNested.String() || schemaType == schemaTypeSetNested.String() || schemaType == schemaTypeMapNested.String()
}

// IsNestedOrArray returns true if the schemaType is nested or an array.
func IsNestedOrArray(schemaType string) bool {
	return IsNested(schemaType) || IsArray(schemaType)
}

// IsList returns true if the schemaType is a list.
func IsList(schemaType string) bool {
	return schemaType == schemaTypeList.String() || schemaType == schemaTypeListNested.String()
}

// IsSet returns true if the schemaType is a set.
func IsSet(schemaType string) bool {
	return schemaType == schemaTypeSet.String()
}

// IsMap returns true if the schemaType is a map.
func IsMap(schemaType string) bool {
	return schemaType == schemaTypeMap.String() || schemaType == schemaTypeMapNested.String()
}

// IsSingle returns true if the schemaType is a singleNestedAttribute.
func IsSingle(schemaType string) bool {
	return schemaType == schemaTypeSingleNested.String()
}
