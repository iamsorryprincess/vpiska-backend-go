package app

import (
	_ "github.com/iamsorryprincess/vpiska-backend-go/docs"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http"
	logging "github.com/iamsorryprincess/vpiska-backend-go/internal/logger"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/security"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/server"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

// @title           Swagger UI
// @version         1.0
// @description     API vpiska.ru
// @BasePath  /api

func Run() {
	logger := logging.NewLogger()
	configuration, configError := parseConfig()

	if configError != nil {
		logger.LogError(configError)
		return
	}

	repositories, repoErr := repository.NewRepositories(configuration.Database.ConnectionString, configuration.Database.DbName)

	if repoErr != nil {
		logger.LogError(repoErr)
		return
	}

	jwtTokenManager := auth.NewJwtManager()
	passwordManager := security.NewPasswordManager()
	services := service.NewServices(repositories, passwordManager, jwtTokenManager)
	handler := http.NewHandler(services, logger, configuration.Server.Port)
	httpServer := server.NewServer(configuration.Server.Port, handler)
	logger.LogError(httpServer.Run())
}
