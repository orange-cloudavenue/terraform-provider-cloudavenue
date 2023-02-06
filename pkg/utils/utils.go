package utils

import (
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Generate a string from a sorted list of strings
func generateStringFromList(list []string) string {
	sort.Strings(list)
	return strings.Join(list, "")
}

// Generate unique UUID from a string
func generateUUID(str string) string {
	return uuid.NewSHA1(uuid.Nil, []byte(str)).String()
}

// Generate unique UUID from a sorted list of strings
func GenerateUUIDFromList(list []string) types.String {
	return types.StringValue(generateUUID(generateStringFromList(list)))
}

// Generate unique UUID from string
func GenerateUUIDFromString(str string) types.String {
	return types.StringValue(generateUUID(str))
}
