package json

import (
	"errors"
	"image-service/internal/model"
	"strings"

	"github.com/gin-gonic/gin"
)

// ConstructJsonResponseSuccess generates a JSON success response using model.SuccessResponse.
func ConstructJsonResponseSuccess(c *gin.Context, data any, optionalParams ...any) {
	statusCode := 200
	message := ""

	if len(optionalParams) > 0 {
		switch v := optionalParams[0].(type) {
		case int:
			statusCode = v
		case string:
			message = v
		}
	}
	if len(optionalParams) > 1 {
		switch v := optionalParams[1].(type) {
		case int:
			statusCode = v
		case string:
			message = v
		}
	}

	resp := model.SuccessResponse{
		Status:  statusCode,
		Data:    data,
		Message: message,
	}

	c.JSON(statusCode, resp)
}

// ConstructJsonResponseError generates a JSON error response using model.ErrorDetails.
func ConstructJsonResponseError(c *gin.Context, err error, statusCode int) {
	details := model.ErrorDetails{
		Status:  statusCode,
		Message: ExtractErrorMessage(err),
	}
	resp := model.ErrorResponse{
		Error: details,
	}

	c.JSON(statusCode, resp)
}

// ExtractErrorMessage tries to extract a cleaner message from a wrapped error.
func ExtractErrorMessage(err error) string {
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		return err.Error()
	}

	parts := strings.SplitN(err.Error(), ": ", 2)
	if len(parts) == 2 {
		return parts[1]
	}

	return err.Error()
}
