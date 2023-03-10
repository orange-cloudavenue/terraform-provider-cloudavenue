// Package mutex provides a simple key/value store for arbitrary mutexes.
package mutex

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// KV is a simple key/value store for arbitrary mutexes. It can be used to
// serialize changes across arbitrary collaborators that share knowledge of the
// keys they must serialize on.
//
// The initial use case is to let aws_security_group_rule resources serialize
// their access to individual security groups based on SG ID.
type KV struct {
	lock  sync.Mutex
	store map[string]*sync.Mutex
}

// NewKV is an implementation of KV.
func NewKV() *KV {
	return &KV{
		store: make(map[string]*sync.Mutex),
	}
}

// KvLock locks the mutex for the given key. Caller is responsible for calling kvUnlock
// for the same key.
func (m *KV) KvLock(ctx context.Context, key string) {
	tflog.Debug(ctx, fmt.Sprintf("Locking %q", key))
	m.get(key).Lock()
	tflog.Debug(ctx, fmt.Sprintf("Locked %q", key))
}

// KvUnlock unlocks the mutex for the given key. Caller must have called kvLock for the same key first.
func (m *KV) KvUnlock(ctx context.Context, key string) {
	tflog.Debug(ctx, fmt.Sprintf("Unlocking %q", key))
	m.get(key).Unlock()
	tflog.Debug(ctx, fmt.Sprintf("Unlocked %q", key))
}

// Returns a mutex for the given key, no guarantee of its lock status.
func (m *KV) get(key string) *sync.Mutex {
	m.lock.Lock()
	defer m.lock.Unlock()
	mutex, ok := m.store[key]
	if !ok {
		mutex = &sync.Mutex{}
		m.store[key] = mutex
	}
	return mutex
}
