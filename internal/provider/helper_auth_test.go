package provider

import (
	"context"
	"testing"
	"time"

	apiclient "github.com/orange-cloudavenue/cloudavenue-sdk-go"
)

func TestGetAuthContextWithTO(t *testing.T) {
	name := "testGetAuthContextWithTO"
	want := "token"

	apiCtx := context.WithValue(context.Background(), apiclient.ContextAccessToken, "token")
	ctx := context.Background()

	tfCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	newCtx, _ := getAuthContextWithTO(apiCtx, tfCtx)
	newToken, _ := newCtx.Value(apiclient.ContextAccessToken).(string)
	if newToken != want {
		t.Errorf("%s: got %s, want %s", name, newToken, want)
	}
}
