// Package helpers provides auth helpers for the CloudAvenue Terraform Provider.
package helpers

import (
	"context"
	"errors"

	apiclient "github.com/orange-cloudavenue/infrapi-sdk-go"
)

// GetAuthContextWithTO is a helper function to create the auth context with the token api and the terraform context with timeout.
func GetAuthContextWithTO(apiCtx, tfCtx context.Context) (context.Context, error) {
	token, tokenIsAString := apiCtx.Value(apiclient.ContextAccessToken).(string)
	if !tokenIsAString {
		return nil, errors.New("token is not a string")
	}
	auth := context.WithValue(tfCtx, apiclient.ContextAccessToken, token)
	return auth, nil
}
