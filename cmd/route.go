package cmd

import "github.com/gin-gonic/gin"

func route(r *gin.Engine, dependency Dependency) {
	r.GET("/health", dependency.HealthcheckAPI.HealthcheckHandlerHTTP)

	userV1 := r.Group("/user/v1")
	userV1.POST("/register", dependency.RegisterAPI.Register)
	userV1.POST("/login", dependency.LoginAPI.Login)
	userV1.DELETE("/logout", dependency.MiddlewareValidateAuth, dependency.LogoutAPI.Logout)
	userV1.PUT("/refresh-token", dependency.MiddlewareRefreshToken, dependency.RefreshTokenAPI.RefreshToken)
}
