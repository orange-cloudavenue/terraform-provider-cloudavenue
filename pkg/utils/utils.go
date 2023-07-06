// Package utils provides utils for Terraform Provider.
package utils

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type tfValuesForUUID interface {
	string | []string
}

// generateUUID generates a unique UUID based on a string.
// This is used to generate a unique ID for a resource.
func generateUUID(str string) string {
	return uuid.NewSHA1(uuid.Nil, []byte(str)).String()
}

// GenerateUUID generates a unique UUID. The value can be a string or a slice of strings.
// This is used to generate a unique ID for a resource.
func GenerateUUID[T tfValuesForUUID](values ...T) types.String {
	s := []string{}

	for _, v := range values {
		s = append(s, fmt.Sprintf("%s", v))
	}

	sort.Strings(s)
	str := strings.Join(s, "")

	return types.StringValue(generateUUID(str))
}

// TakeBoolPointer accepts a boolean and returns a pointer to this value.
func TakeBoolPointer(value bool) *bool {
	return &value
}

// TakeIntPointer accepts an int and returns a pointer to this value.
func TakeIntPointer(x int) *int {
	return &x
}

// TakeInt64Pointer accepts an int64 and returns a pointer to this value.
func TakeInt64Pointer(x int64) *int64 {
	return &x
}

// StringValueOrNull return a null StringValue if value is "" or return StringValue(value) if not.
func StringValueOrNull(value string) basetypes.StringValue {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

// SortMapStringByKeys sorts a map[string]string by keys.
func SortMapStringByKeys[T any](m map[string]T) map[string]T {
	sortedKeys := make([]string, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	sortedMap := make(map[string]T)
	for _, k := range sortedKeys {
		sortedMap[k] = m[k]
	}

	return sortedMap
}

// SliceTypesStringToSliceString converts a slice of types.String to a slice of string.
func SliceTypesStringToSliceString(slice []types.String) []string {
	var result []string
	for _, s := range slice {
		result = append(result, s.ValueString())
	}
	return result
}

type OpenAPIValues []string

// OpenApiReferenceToSliceID converts a slice of OpenApiReference to a slice of ID.
func OpenAPIReferenceToSliceID(slice []govcdtypes.OpenApiReference) OpenAPIValues {
	var result OpenAPIValues
	for _, s := range slice {
		result = append(result, s.ID)
	}
	return result
}

// OpenAPIReferenceToSliceName converts a slice of OpenApiReference to a slice of Name.
func OpenAPIReferenceToSliceName(slice []govcdtypes.OpenApiReference) OpenAPIValues {
	var result OpenAPIValues
	for _, s := range slice {
		result = append(result, s.Name)
	}
	return result
}

// ToTerraformTypes converts a slice of string to a slice of types.String.
func (o *OpenAPIValues) ToTerraformTypesString() []types.String {
	var result []types.String
	for _, s := range *o {
		result = append(result, types.StringValue(s))
	}
	return result
}

// ToTerraformTypesSet converts a slice of string to a slice of types.StringSet.
func (o OpenAPIValues) ToTerraformTypesStringSet(ctx context.Context) basetypes.SetValue {
	x, _ := types.SetValueFrom(ctx, types.StringType, o)
	return x
}
