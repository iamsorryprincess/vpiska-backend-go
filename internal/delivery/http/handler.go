package http

import (
	"net/http"

	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(services *service.Services, logger logger.Logger, tokenManager auth.TokenManager) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	swaggerHandler := httpSwagger.Handler(httpSwagger.URL("doc.json"))
	mux.Handle("/swagger/", swaggerHandler)
	handler := v1.NewHandler(logger, services, tokenManager)
	handler.InitAPI(mux)
	websocket.InitWebsocketsRoutes(mux, logger, tokenManager, services.Events, services.Publisher)
	return handler.Recover(mux)
}
