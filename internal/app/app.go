package app

import (
	_ "github.com/iamsorryprincess/vpiska-backend-go/docs"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/server"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logging"
)

// @title           Swagger UI
// @version         1.0
// @description     API vpiska.ru
// @BasePath  /api

func Run() {
	logger := logging.NewLogger()
	configuration, err := parseConfig()

	if err != nil {
		logger.LogError(err)
		return
	}

	repositories, err := repository.NewRepositories(configuration.Database.ConnectionString, configuration.Database.DbName)

	if err != nil {
		logger.LogError(err)
		return
	}

	jwtTokenManager := auth.NewJwtManager()
	passwordManager := hash.NewPasswordHashManager()
	services := service.NewServices(repositories, passwordManager, jwtTokenManager)
	handler := http.NewHandler(services, logger, configuration.Server.Port)
	httpServer := server.NewServer(configuration.Server.Port, handler)
	logger.LogError(httpServer.Run())
}
