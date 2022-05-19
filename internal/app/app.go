package app

import (
	"context"
	"log"
	"os"
	"time"

	_ "github.com/iamsorryprincess/vpiska-backend-go/docs"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/server"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/storage"
)

// @title           Swagger UI
// @version         1.0
// @description     API vpiska.ru
// @BasePath  /api

// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization

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

	jwtDuration := time.Hour * 24 * 3
	jwtTokenManager := auth.NewJwtManager(configuration.JWT.Key, configuration.JWT.Issuer, configuration.JWT.Audience, jwtDuration)
	passwordManager := hash.NewPasswordHashManager()

	mediasMetadata, err := repositories.Media.GetAll(context.Background())

	if err != nil {
		errorLogger.Println(err)
		return
	}

	mediaIds := make([]string, len(mediasMetadata))

	for index, mediaMetadata := range mediasMetadata {
		mediaIds[index] = mediaMetadata.ID
	}

	fileStorage, err := storage.NewLocalFileStorage("media", mediaIds)

	if err != nil {
		errorLogger.Println(err)
		return
	}

	services, err := service.NewServices(repositories, passwordManager, jwtTokenManager, fileStorage)

	if err != nil {
		errorLogger.Println(err)
		return
	}

	handler := http.NewHandler(services, errorLogger, jwtTokenManager, configuration.Server.Port)
	httpServer := server.NewServer(configuration.Server.Port, handler)
	errorLogger.Println(httpServer.Run())
}
