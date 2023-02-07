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
