package custom_errors

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	base := ErrBadRequest
	message := "invalid image format"
	err := NewError(base, message)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, base))
	assert.Contains(t, err.Error(), message)
}

func TestMapErrors_KnownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	wrapped := NewError(ErrNotFound, "album does not exist")

	MapErrors(c, wrapped)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "album does not exist")
	assert.JSONEq(t, `{
		"error": {
			"status": 404,
			"message": "album does not exist"
		}
	}`, w.Body.String())
}

func TestMapErrors_UnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	err := errors.New("some unknown error")

	MapErrors(c, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "unexpected error")
	assert.JSONEq(t, `{
		"error": {
			"status": 500,
			"message": "unexpected error"
		}
	}`, w.Body.String())
}
