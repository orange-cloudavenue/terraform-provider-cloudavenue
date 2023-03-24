package cloudavenue

import (
	"context"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/provider/common/mutex"
)

var cAMutexKV = mutex.NewKV()

const keyLock = "cloudavenue:customer:api:lock"

// Lock
// lock call to the cloudavenue customer API.
func Lock(ctx context.Context) {
	cAMutexKV.KvLock(ctx, keyLock)
}

// Unlock
// unlock call to the cloudavenue customer API.
func Unlock(ctx context.Context) {
	cAMutexKV.KvUnlock(ctx, keyLock)
}
