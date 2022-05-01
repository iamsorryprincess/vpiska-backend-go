package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/database"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/identity"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func Run() {
	configuration, configError := parseConfig()

	if configError != nil {
		log.Fatal(configError)
		return
	}

	userRepository, userRepoErr := database.InitUserRepository(
		configuration.Database.ConnectionString,
		configuration.Database.DbName,
		"users")

	if userRepoErr != nil {
		log.Fatal(userRepoErr)
		return
	}

	securityProvider := identity.InitPasswordHashProvider()
	identityProvider := identity.InitJwtTokenProvider()
	userService := service.InitUserService(userRepository, securityProvider, identityProvider)

	http.Handle("/api/v1/users/create", v1.CreateUserHandler(userService))
	http.Handle("/api/v1/users/login", v1.LoginUserHandler(userService))

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", configuration.Server.Host, configuration.Server.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
