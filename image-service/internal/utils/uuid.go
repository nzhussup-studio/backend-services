package utils

import (
	"log"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	return id.String()
}
