package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRouterReturnsGinEngine(t *testing.T) {
	r := SetupRouter()
	if r == nil {
		t.Fatal("SetupRouter() returned nil")
	}
	_, ok := interface{}(r).(*gin.Engine)
	if !ok {
		t.Fatalf("SetupRouter() did not return *gin.Engine, got %T", r)
	}
}
