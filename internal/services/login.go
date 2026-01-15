package services

import (
	"context"
	"fmt"
	"time"

	"ewallet-ums/helpers"
	"ewallet-ums/internal/interfaces"
	"ewallet-ums/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	UserRepo interfaces.IUserRepository
}

func (s *LoginService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	resp := models.LoginResponse{}
	now := time.Now()

	userDetail, err := s.UserRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return resp, fmt.Errorf("failed to get user by username, %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(req.Password)); err != nil {
		return resp, fmt.Errorf("incorrect password, %v", err)
	}

	token, err := helpers.GenerateToken(ctx, userDetail.ID, userDetail.Username, userDetail.FullName, "token", userDetail.Email, now)
	if err != nil {
		return resp, fmt.Errorf("failed to generate token, %v", err)
	}

	refreshToken, err := helpers.GenerateToken(ctx, userDetail.ID, userDetail.Username, userDetail.FullName, "refresh_token", userDetail.Email, now)
	if err != nil {
		return resp, fmt.Errorf("failed to generate refresh token, %v", err)
	}

	userSession := &models.UserSession{
		UserID:              userDetail.ID,
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(helpers.MapTypeToken["token"]),
		RefreshTokenExpired: now.Add(helpers.MapTypeToken["refresh_token"]),
	}
	err = s.UserRepo.InsertNewUserSession(ctx, userSession)
	if err != nil {
		return resp, fmt.Errorf("failed to insert new session, %v", err)
	}

	resp.UserID = userDetail.ID
	resp.Username = userDetail.Username
	resp.FullName = userDetail.FullName
	resp.Email = userDetail.Email
	resp.Token = token
	resp.RefreshToken = refreshToken

	return resp, nil
}
