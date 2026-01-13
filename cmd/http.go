package cmd

import (
	"log"

	"ewallet-ums/helpers"
	"ewallet-ums/internal/api"
	"ewallet-ums/internal/repository"
	"ewallet-ums/internal/services"

	"github.com/gin-gonic/gin"
)

func ServeHTTP() {
	healthcheckSvc := &services.Healthcheck{}
	healthcheckAPI := &api.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	registerRepo := &repository.RegisterRepository{
		DB: helpers.DB,
	}
	registerSvc := &services.RegisterService{
		RegisterRepo: registerRepo,
	}
	registerAPI := &api.RegisterHandler{
		RegisterService: registerSvc,
	}

	r := gin.Default()
	r.GET("/health", healthcheckAPI.HealthcheckHandlerHTTP)

	userV1 := r.Group("/users/v1")
	userV1.POST("/register", registerAPI.Register)

	err := r.Run(":" + helpers.GetEnv("PORT", "8080"))
	if err != nil {
		log.Fatal(err)
	}
}
