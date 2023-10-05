// Package helpers provides auth helpers for the CloudAvenue Terraform Provider.
package helpers

import (
	"context"
	"testing"
	"time"

	apiclient "github.com/orange-cloudavenue/infrapi-sdk-go"
)

func TestGetAuthContextWithTO(t *testing.T) {
	name := "testGetAuthContextWithTO"

	apiCtx := context.WithValue(context.Background(), apiclient.ContextAccessToken, "token")
	ctx := context.Background()

	tfCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	newCtx, _ := GetAuthContextWithTO(apiCtx, tfCtx)
	newToken, _ := newCtx.Value(apiclient.ContextAccessToken).(string)
	if want := "token"; newToken != want {
		t.Errorf("%s: got %s, want %s", name, newToken, want)
	}
}
