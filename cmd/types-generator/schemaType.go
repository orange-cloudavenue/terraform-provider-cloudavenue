package main

type schemaType string

const (
	schemaTypeBool       schemaType = "schema.BoolAttribute"
	schemaTypeInt64      schemaType = "schema.Int64Attribute"
	schemaTypeFloat64    schemaType = "schema.Float64Attribute"
	schemaTypeString     schemaType = "schema.StringAttribute"
	schemaTypeList       schemaType = "schema.ListAttribute"
	schemaTypeSet        schemaType = "schema.SetAttribute"
	schemaTypeMap        schemaType = "schema.MapAttribute"
	schemaTypeListNested schemaType = "schema.ListNestedAttribute"
	schemaTypeSetNested  schemaType = "schema.SetNestedAttribute"
	schemaTypeMapNested  schemaType = "schema.MapNestedAttribute"
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
	}

	return ""
}

// String returns the string representation of the schemaType.
func (s schemaType) String() string {
	return string(s)
}

// ToTerraformType returns the terraform type of the schemaType.
func (s schemaType) ToTerraformType() string {
	switch s {
	case schemaTypeBool:
		return "types.Bool"
	case schemaTypeInt64:
		return "types.Int64"
	case schemaTypeFloat64:
		return "types.Float64"
	case schemaTypeString:
		return "types.String"
	case schemaTypeList, schemaTypeListNested:
		return "types.List"
	case schemaTypeSet, schemaTypeSetNested:
		return "types.Set"
	case schemaTypeMap, schemaTypeMapNested:
		return "types.Map"
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
		return "types.BoolNull"
	case schemaTypeInt64:
		return "types.Int64Null"
	case schemaTypeFloat64:
		return "types.Float64Null"
	case schemaTypeString:
		return "types.StringNull"
	case schemaTypeList, schemaTypeListNested:
		return "types.ListNull"
	case schemaTypeSet, schemaTypeSetNested:
		return "types.SetNull"
	case schemaTypeMap, schemaTypeMapNested:
		return "types.MapNull"
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

// IsNested returns true if the schemaType is nested.
func IsNested(schemaType string) bool {
	switch schemaType {
	case schemaTypeListNested.String(), schemaTypeSetNested.String(), schemaTypeMapNested.String():
		return true
	}

	return false
}

// IsList returns true if the schemaType is a list.
func IsList(schemaType string) bool {
	switch schemaType {
	case schemaTypeList.String(), schemaTypeListNested.String():
		return true
	}

	return false
}

// IsSet returns true if the schemaType is a set.
func IsSet(schemaType string) bool {
	switch schemaType {
	case schemaTypeSet.String(), schemaTypeSetNested.String():
		return true
	}

	return false
}

// IsMap returns true if the schemaType is a map.
func IsMap(schemaType string) bool {
	switch schemaType {
	case schemaTypeMap.String(), schemaTypeMapNested.String():
		return true
	}

	return false
}
