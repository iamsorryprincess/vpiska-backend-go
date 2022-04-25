package app

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/commands"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/infrastructure/database"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/infrastructure/identity"
)

func Run() {
	userRepository := database.InitUserRepository()
	passwordHashProvider := identity.InitPasswordHashProvider()
	jwtIdentityProvider := identity.InitJwtIdentityProvider()
	createUserHandler := commands.InitCreateUserHandler(userRepository, passwordHashProvider, jwtIdentityProvider)
	usersHandler := v1.InitUsersHandler(createUserHandler)
	handler := v1.InitHandler(usersHandler)

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	v1.InitV1ApiHandlers(r, handler)
	err := http.ListenAndServe("localhost:5000", r)
	if err != nil {
		log.Fatal(err)
	}
}
