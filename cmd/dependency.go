package cmd

import (
	"ewallet-ums/helpers"
	"ewallet-ums/internal/api"
	"ewallet-ums/internal/interfaces"
	"ewallet-ums/internal/repository"
	"ewallet-ums/internal/services"
)

type Dependency struct {
	UserRepo interfaces.IUserRepository

	HealthcheckAPI  interfaces.IHealthcheckHandler
	RegisterAPI     interfaces.IRegisterHandler
	LoginAPI        interfaces.ILoginHandler
	LogoutAPI       interfaces.ILogoutHandler
	RefreshTokenAPI interfaces.IRefreshTokenHandler

	TokenValidationAPI *api.TokenValidationHandler
}

func dependencyInject() Dependency {
	healthcheckSvc := &services.Healthcheck{}
	healthcheckAPI := &api.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	userRepo := &repository.UserRepository{
		DB: helpers.DB,
	}

	registerSvc := &services.RegisterService{
		UserRepo: userRepo,
	}

	registerAPI := &api.RegisterHandler{
		RegisterService: registerSvc,
	}

	loginSvc := &services.LoginService{
		UserRepo: userRepo,
	}

	loginAPI := &api.LoginHandler{
		LoginService: loginSvc,
	}

	logoutSvc := &services.LogoutService{
		UserRepo: userRepo,
	}

	logoutAPI := &api.LogoutHandler{
		LogoutService: logoutSvc,
	}

	refreshTokenSvc := &services.RefreshTokenService{
		UserRepo: userRepo,
	}

	refreshTokenAPI := &api.RefreshTokenHandler{
		RefreshTokenService: refreshTokenSvc,
	}

	tokenValidationSvc := &services.TokenValidationService{
		UserRepo: userRepo,
	}

	tokenValidationAPI := &api.TokenValidationHandler{
		TokenValidationService: tokenValidationSvc,
	}

	return Dependency{
		UserRepo:           userRepo,
		HealthcheckAPI:     healthcheckAPI,
		RegisterAPI:        registerAPI,
		LoginAPI:           loginAPI,
		LogoutAPI:          logoutAPI,
		RefreshTokenAPI:    refreshTokenAPI,
		TokenValidationAPI: tokenValidationAPI,
	}
}
