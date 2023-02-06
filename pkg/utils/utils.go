package utils

import (
	"errors"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// generateUUID generates a unique UUID based on a string.
// This is used to generate a unique ID for a resource.
func generateUUID(str string) string {
	return uuid.NewSHA1(uuid.Nil, []byte(str)).String()
}

// GenerateUUID generates a unique UUID. The value can be a string or a slice of strings.
// If the value is a slice of strings, the strings will be sorted before generating the UUID.
// This is used to generate a unique ID for a resource.
func GenerateUUID(value interface{}) (types.String, error) {
	var val string

	switch v := value.(type) {
	case []string:
		sort.Strings(v)
		val = strings.Join(v, "")
	case string:
		val = v
	default:
		return types.StringNull(), errors.New("invalid type")
	}

	return types.StringValue(generateUUID(val)), nil
}
