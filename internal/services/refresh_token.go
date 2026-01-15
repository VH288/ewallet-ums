package services

import (
	"context"
	"fmt"
	"time"

	"ewallet-ums/helpers"
	"ewallet-ums/internal/interfaces"
	"ewallet-ums/internal/models"
)

type RefreshTokenService struct {
	UserRepo interfaces.IUserRepository
}

func (s *RefreshTokenService) RefreshToken(ctx context.Context, refreshToken string, tokenClaim helpers.ClaimToken) (models.RefreshTokenResponse, error) {
	resp := models.RefreshTokenResponse{}

	token, err := helpers.GenerateToken(ctx, tokenClaim.UserID, tokenClaim.Username, tokenClaim.FullName, "refresh_token", tokenClaim.Email, time.Now())
	if err != nil {
		return resp, fmt.Errorf("failed to generate new token %v", err)
	}

	err = s.UserRepo.UpdateTokenByRefreshToken(ctx, token, refreshToken)
	if err != nil {
		return resp, fmt.Errorf("failed to update new token %v", err)
	}

	resp.Token = token
	return resp, nil
}
