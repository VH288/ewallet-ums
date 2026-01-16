package api

import (
	"context"
	"fmt"

	"ewallet-ums/cmd/proto/tokenvalidation"
	"ewallet-ums/constants"
	"ewallet-ums/helpers"
	"ewallet-ums/internal/interfaces"
)

type TokenValidationHandler struct {
	tokenvalidation.UnimplementedTokenValidationServer
	TokenValidationService interfaces.ITokenValidationService
}

func (s *TokenValidationHandler) ValidateToken(ctx context.Context, req *tokenvalidation.TokenRequest) (*tokenvalidation.TokenResponse, error) {
	var (
		token = req.GetToken()
		log   = helpers.Logger
	)

	if token == "" {
		err := fmt.Errorf("token is empty")
		log.Error(err)
		return &tokenvalidation.TokenResponse{
			Message: err.Error(),
		}, nil
	}

	claimToken, err := s.TokenValidationService.TokenValidation(ctx, token)
	if err != nil {
		log.Error(err)
		return &tokenvalidation.TokenResponse{
			Message: err.Error(),
		}, nil
	}

	fmt.Println("user id: ", claimToken.UserID)

	return &tokenvalidation.TokenResponse{
		Message: constants.SuccessMessage,
		Data: &tokenvalidation.UserData{
			UserId:   int64(claimToken.UserID),
			Username: claimToken.Username,
			FullName: claimToken.FullName,
		},
	}, nil
}
