package json

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestConstructJsonResponseSuccess_Defaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := gin.H{"key": "value"}

	ConstructJsonResponseSuccess(c, data)

	assert.Equal(t, http.StatusOK, w.Code)
	expected := `{"status":200,"data":{"key":"value"}}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestConstructJsonResponseSuccess_WithStatusCodeAndMessage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := gin.H{"foo": "bar"}

	ConstructJsonResponseSuccess(c, data, 201, "created")

	assert.Equal(t, 201, w.Code)
	expected := `{"status":201,"data":{"foo":"bar"},"message":"created"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestConstructJsonResponseSuccess_WithMessageFirst(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := gin.H{"id": 1}

	// Message string first, status code second
	ConstructJsonResponseSuccess(c, data, "ok", 202)

	assert.Equal(t, 202, w.Code)
	expected := `{"status":202,"data":{"id":1},"message":"ok"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestConstructJsonResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	err := errors.New("something went wrong")

	ConstructJsonResponseError(c, err, http.StatusBadRequest)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	expected := `{"error":{"status":400,"message":"something went wrong"}}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestExtractErrorMessage_Simple(t *testing.T) {
	err := errors.New("basic error message")
	msg := ExtractErrorMessage(err)
	assert.Equal(t, "basic error message", msg)
}

func TestExtractErrorMessage_Wrapped(t *testing.T) {
	baseErr := errors.New("original error")
	wrapped := errors.Join(errors.New("ignored"), baseErr)
	msg := ExtractErrorMessage(wrapped)
	// Fallback to full message since .Unwrap() doesn't behave like errors.Join
	assert.Equal(t, wrapped.Error(), msg)
}
