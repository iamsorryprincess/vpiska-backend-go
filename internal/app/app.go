package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/database"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1/handler"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/identity"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

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
	userHandler := handler.NewUserHandler(userService)

	http.HandleFunc("/api/v1/users/create", userHandler.CreateUser)
	http.HandleFunc("/api/v1/users/login", userHandler.LoginUser)
	http.HandleFunc("/api/v1/users/password/change", userHandler.ChangePassword)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", configuration.Server.Host, configuration.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
