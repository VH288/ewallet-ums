package services

import (
	"context"
	"fmt"

	"ewallet-ums/helpers"
	"ewallet-ums/internal/interfaces"
)

type TokenValidationService struct {
	UserRepo interfaces.IUserRepository
}

func (s *TokenValidationService) TokenValidation(ctx context.Context, token string) (*helpers.ClaimToken, error) {
	var claimToken *helpers.ClaimToken

	claimToken, err := helpers.ValidateToken(ctx, token)
	if err != nil {
		return claimToken, fmt.Errorf("failed to validate token: %v", err)
	}

	_, err = s.UserRepo.GetUserSessionByToken(ctx, token)
	if err != nil {
		return claimToken, fmt.Errorf("failed to get user session: %v", err)
	}

	return claimToken, nil
}
