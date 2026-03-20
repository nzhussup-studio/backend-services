package controller

import (
	"net/http"

	customerrors "account-service/internal/errors"
	"account-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AccountController struct {
	service *service.AccountService
}

func NewAccountController(service *service.AccountService) *AccountController {
	return &AccountController{service: service}
}

func (ac *AccountController) DeleteCurrentAccount(c *gin.Context) {
	subject, ok := c.Get("subject")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := ac.service.DeleteCurrentAccount(c, subject.(string)); err != nil {
		status := customerrors.StatusCode(err)
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
