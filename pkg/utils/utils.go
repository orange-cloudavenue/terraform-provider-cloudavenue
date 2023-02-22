// Package utils provides utils for Terraform Provider.
package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
