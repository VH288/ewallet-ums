package interfaces

import (
	"context"

	"ewallet-ums/internal/models"
)

type IRegisterService interface {
	Register(ctx context.Context, request *models.User) (any, error)
}
