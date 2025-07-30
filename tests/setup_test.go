//go:build integration

package tests

import (
	"testing"

	"github.com/denkhaus/tensorzero"
)

// setupTestClient creates a test client for integration tests
func setupTestClient(t testing.TB) tensorzero.Gateway {
	client := tensorzero.NewHTTPGateway("http://localhost:3000")
	t.Cleanup(func() {
		client.Close()
	})
	return client
}