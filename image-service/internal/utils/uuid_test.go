package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUID(t *testing.T) {
	id := GenerateUUID()
	if id == "" {
		t.Error("Expected a non-empty UUID string")
	}

	_, err := uuid.Parse(id)
	if err != nil {
		t.Errorf("Generated UUID is not valid: %v", err)
	}
}
