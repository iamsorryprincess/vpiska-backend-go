package app

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/iamsorryprincess/vpiska-backend-go/docs"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/database"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/identity"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Swagger UI
// @version         1.0
// @description     API vpiska.ru
// @BasePath  /api/v1

func Run() {
	configuration, configError := parseConfig()

	if configError != nil {
		log.Fatal(configError)
		return
	}

	userRepository, userRepoErr := database.NewUserRepository(
		configuration.Database.ConnectionString,
		configuration.Database.DbName,
		"users")

	if userRepoErr != nil {
		log.Fatal(userRepoErr)
		return
	}

	securityProvider := identity.NewPasswordHashProvider()
	identityProvider := identity.NewJwtTokenProvider()
	userService := service.NewUserService(userRepository, securityProvider, identityProvider)
	userHandler := v1.NewUserHandler(userService)

	swaggerUrl := fmt.Sprintf("http://localhost:%d/swagger/doc.json", configuration.Server.Port)
	http.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL(swaggerUrl)))
	http.HandleFunc("/api/v1/users/create", userHandler.CreateUser)
	http.HandleFunc("/api/v1/users/login", userHandler.LoginUser)
	http.HandleFunc("/api/v1/users/password/change", userHandler.ChangePassword)

	err := http.ListenAndServe(fmt.Sprintf(":%d", configuration.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
