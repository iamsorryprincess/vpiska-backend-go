package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/iamsorryprincess/vpiska-backend-go/docs"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/config"
	appHttp "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/repository"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/server"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/hash"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
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
	appLogger, logFile, err := logger.NewLogLogger()
	defer logFile.Close()

	if err != nil {
		log.Fatal(err)
		return
	}

	configuration, err := config.Parse()

	if err != nil {
		appLogger.LogError(err)
		return
	}

	repositories, err := repository.NewRepositories(configuration.Database.ConnectionString, configuration.Database.DbName)

	if err != nil {
		appLogger.LogError(err)
		return
	}

	jwtDuration := time.Hour * 24 * 3
	jwtTokenManager := auth.NewJwtManager(configuration.JWT.Key, configuration.JWT.Issuer, configuration.JWT.Audience, jwtDuration)
	passwordManager, err := hash.NewPasswordHashManager(configuration.Hash.Key)

	if err != nil {
		appLogger.LogError(err)
		return
	}

	fileStorage, err := storage.NewLocalFileStorage("media")

	if err != nil {
		appLogger.LogError(err)
		return
	}

	services, err := service.NewServices(appLogger, repositories, passwordManager, jwtTokenManager, fileStorage)

	if err != nil {
		appLogger.LogError(err)
		return
	}

	handler := appHttp.NewHandler(services, appLogger, jwtTokenManager, configuration.Server.Port)
	httpServer := server.NewServer(configuration.Server.Port, handler)

	go func() {
		if err = httpServer.Run(); err != nil && errors.Is(err, http.ErrServerClosed) {
			appLogger.LogError(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	services.Publisher.CloseAll()
	appLogger.LogInfo("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err = httpServer.Stop(ctx); err != nil {
		appLogger.LogError(err)
		return
	}

	appLogger.LogInfo("Server exiting")
}
