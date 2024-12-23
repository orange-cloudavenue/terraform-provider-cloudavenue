package common

import (
	"context"
	"regexp"

	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
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

func FromOpenAPIReferenceID(_ context.Context, apiRefs []govcdtypes.OpenApiReference) (values []string) {
	if len(apiRefs) == 0 {
		return nil
	}

	values = make([]string, 0)
	for _, apiRef := range apiRefs {
		values = append(values, apiRef.ID)
	}

	return
}

func ToOpenAPIReferenceID(ctx context.Context, attribute supertypes.SetValueOf[string]) (apiRefs []govcdtypes.OpenApiReference, diags diag.Diagnostics) {
	if attribute.IsKnown() {
		values, d := attribute.Get(ctx)
		if d.HasError() {
			diags.Append(d...)
			return
		}

		if len(values) == 0 {
			return nil, nil
		}

		openAPIReferences := make([]govcdtypes.OpenApiReference, 0)
		for _, id := range values {
			openAPIReferences = append(openAPIReferences, govcdtypes.OpenApiReference{
				ID: id,
			})
		}
		return openAPIReferences, nil
	}

	return nil, nil
}
