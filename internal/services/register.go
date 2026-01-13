package services

import (
	"context"

	"ewallet-ums/internal/interfaces"
	"ewallet-ums/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterService struct {
	RegisterRepo interfaces.IRegisterRepository
}

func (s *RegisterService) Register(ctx context.Context, request *models.User) (any, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	request.Password = string(hashPassword)

	err = s.RegisterRepo.InsertNewUser(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := request
	resp.Password = ""
	return resp, nil
}
