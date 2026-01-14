package interfaces

import (
	"context"

	"ewallet-ums/internal/models"

	"github.com/gin-gonic/gin"
)

type IRegisterService interface {
	Register(ctx context.Context, request *models.User) (any, error)
}

type IRegisterHandler interface {
	Register(c *gin.Context)
}
