package common

import (
	"regexp"
)

// ExtractUUID finds an UUID in the input string
// Returns an empty string if no UUID was found.
func ExtractUUID(input string) string {
	reGetID := regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})`)
	matchListIDs := reGetID.FindAllStringSubmatch(input, -1)
	if len(matchListIDs) > 0 && len(matchListIDs[0]) > 0 {
		return matchListIDs[len(matchListIDs)-1][len(matchListIDs[0])-1]
	}
	return ""
}
