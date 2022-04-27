package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/commands"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/infrastructure/database"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/infrastructure/identity"
)

func Run() {
	configuration, configErr := parseConfig()

	if configErr != nil {
		log.Fatal(configErr)
		return
	}

	userRepository, repoError := database.InitUserRepository(
		configuration.Database.ConnectionString,
		configuration.Database.DbName,
		"users")

	if repoError != nil {
		log.Fatal(repoError)
		return
	}

	passwordHashProvider := identity.InitPasswordHashProvider()
	jwtIdentityProvider := identity.InitJwtIdentityProvider()
	createUserHandler := commands.InitCreateUserHandler(userRepository, passwordHashProvider, jwtIdentityProvider)
	usersHandler := v1.InitUsersHandler(createUserHandler)
	handler := v1.InitHandler(usersHandler)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	v1.InitV1ApiHandlers(r, handler)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", configuration.Server.Host, configuration.Server.Port), r)
	if err != nil {
		log.Fatal(err)
	}
}
