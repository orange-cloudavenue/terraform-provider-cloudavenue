package provider

import (
	"context"
	"errors"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

// getAuthContextWithTO is a helper function to create the auth context with the token api and the terraform context with timeout.
func getAuthContextWithTO(apiCtx, tfCtx context.Context) (context.Context, error) {
	token, tokenIsAString := apiCtx.Value(apiclient.ContextAccessToken).(string)
	if !tokenIsAString {
		return nil, errors.New("token is not a string")
	}
	auth := context.WithValue(tfCtx, apiclient.ContextAccessToken, token)
	return auth, nil
}
