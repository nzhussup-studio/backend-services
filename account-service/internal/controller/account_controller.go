package controller

import (
	"net/http"

	customerrors "account-service/internal/errors"
	"account-service/internal/model"
	"account-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service *service.AccountService
}

func NewAccountController(service *service.AccountService) *AccountController {
	return &AccountController{service: service}
}

// DeleteCurrentAccount godoc
// @Summary Delete the currently authenticated account
// @Description Logs out and deletes the current Keycloak user identified by the bearer token subject.
// @Tags Account
// @Produce json
// @Success 204 {string} string "No Content"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 403 {object} model.ErrorResponse "Forbidden"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /v1/account [delete]
// @Security ApiKeyAuth
func (ac *AccountController) DeleteCurrentAccount(c *gin.Context) {
	subject, ok := c.Get("subject")
	username, usernameOK := c.Get("username")
	email, emailOK := c.Get("email")
	if !ok && !usernameOK && !emailOK {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
		return
	}

	if err := ac.service.DeleteCurrentAccount(
		c,
		valueOrEmpty(subject, ok),
		valueOrEmpty(username, usernameOK),
		valueOrEmpty(email, emailOK),
	); err != nil {
		status := customerrors.StatusCode(err)
		c.JSON(status, model.ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func valueOrEmpty(value any, ok bool) string {
	if !ok {
		return ""
	}

	str, castOK := value.(string)
	if !castOK {
		return ""
	}

	return str
}
