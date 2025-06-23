package http_transport

import (
	"github.com/google/uuid"
	"testing"
)

func TestGenerateJwtKeyPair(t *testing.T) {
	t.Skip()
	// Generate a key pair
	err := generateAndStoreKeys(uuid.NewString())
	if err != nil {
		t.Fatalf("Failed to generate JWT key pair: %v", err)
	}
}
