package services

import (
	"context"

	"ewallet-ums/internal/interfaces"
)

type LogoutService struct {
	UserRepo interfaces.IUserRepository
}

func (s *LogoutService) Logout(ctx context.Context, token string) error {
	return s.UserRepo.DeleteUserSession(ctx, token)
}
