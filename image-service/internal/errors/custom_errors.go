package custom_errors

import (
	"errors"
	"fmt"
	"image-service/internal/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrConflict       = errors.New("conflict")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)

var ErrorMap = map[error]int{
	ErrNotFound:       http.StatusNotFound,
	ErrInternalServer: http.StatusInternalServerError,
	ErrBadRequest:     http.StatusBadRequest,
	ErrConflict:       http.StatusConflict,
	ErrUnauthorized:   http.StatusUnauthorized,
	ErrForbidden:      http.StatusForbidden,
}

func NewError(baseError error, message string) error {
	return fmt.Errorf("%w: %s", baseError, message)
}

func MapErrors(c *gin.Context, err error) {
	mapper := errors.Unwrap(err)
	if mapper != nil {
		statusCode, ok := ErrorMap[mapper]
		if ok {
			json.ConstructJsonResponseError(c, err, statusCode)
			return
		}
	}
	json.ConstructJsonResponseError(c, errors.New("unexpected error"), http.StatusInternalServerError)
}
