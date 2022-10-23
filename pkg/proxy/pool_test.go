package proxy

import "testing"

func TestPool(t *testing.T) {
	pool := NewPool([]string{"localhost:80:user:pass"}, "http")
	_ = pool.GetTransport()
}
