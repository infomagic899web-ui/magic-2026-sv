package middlewares

import (
	"sync"
	"time"
)

type Revoker struct {
	csrf sync.Map // map[string]time.Time
	rsp  sync.Map
}

var TokenRevoker = &Revoker{}

// Mark token as revoked
func (r *Revoker) RevokeCSRF(token string) {
	r.csrf.Store(token, time.Now())
}
func (r *Revoker) RevokeRSP(token string) {
	r.rsp.Store(token, time.Now())
}

// Check if token is revoked
func (r *Revoker) IsCSRFRevoked(token string) bool {
	_, exists := r.csrf.Load(token)
	return exists
}
func (r *Revoker) IsRSPRevoked(token string) bool {
	_, exists := r.rsp.Load(token)
	return exists
}

// Cleanup expired revoked tokens (optional, for memory)
func (r *Revoker) Cleanup() {
	expiry := 10 * time.Second // slightly longer than 5s token lifetime
	now := time.Now()

	r.csrf.Range(func(key, value interface{}) bool {
		if t, ok := value.(time.Time); ok && now.Sub(t) > expiry {
			r.csrf.Delete(key)
		}
		return true
	})
	r.rsp.Range(func(key, value interface{}) bool {
		if t, ok := value.(time.Time); ok && now.Sub(t) > expiry {
			r.rsp.Delete(key)
		}
		return true
	})
}
