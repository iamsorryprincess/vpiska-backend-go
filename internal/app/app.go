package app

import (
	"log"
	"os"

	_ "github.com/iamsorryprincess/vpiska-backend-go/docs"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/server"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
)

// @title           Swagger UI
// @version         1.0
// @description     API vpiska.ru
// @BasePath  /api

func Run() {
	logFile, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	defer logFile.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	errorLogger := log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	configuration, err := parseConfig()

	if err != nil {
		errorLogger.Println(err)
		return
	}

	repositories, err := repository.NewRepositories(configuration.Database.ConnectionString, configuration.Database.DbName)

	if err != nil {
		errorLogger.Println(err)
		return
	}

	jwtTokenManager := auth.NewJwtManager()
	passwordManager := hash.NewPasswordHashManager()

	services, err := service.NewServices(repositories, passwordManager, jwtTokenManager)

	if err != nil {
		errorLogger.Println(err)
		return
	}

	handler := http.NewHandler(services, errorLogger, configuration.Server.Port)
	httpServer := server.NewServer(configuration.Server.Port, handler)
	errorLogger.Println(httpServer.Run())
}
