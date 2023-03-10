package common

import (
	"regexp"
	"strings"
)

// NormalizeID checks if the ID contains a wanted prefix
// If it does, the function returns the original ID.
// Otherwise, it returns the prefix + the ID.
func NormalizeID(prefix, id string) string {
	if strings.Contains(id, prefix) {
		return id
	}
	return prefix + id
}

// ExtractUUID finds an UUID in the input string
// Returns an empty string if no UUID was found.
func ExtractUUID(input string) string {
	reGetID := regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)
	matchListIds := reGetID.FindAllStringSubmatch(input, -1)
	if len(matchListIds) > 0 && len(matchListIds[0]) > 0 {
		return matchListIds[len(matchListIds)-1][len(matchListIds[0])-1]
	}
	return ""
}
