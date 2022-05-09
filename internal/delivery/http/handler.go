package http

import (
	"fmt"
	"net/http"

	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/logger"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(services *service.Services, logger logger.Logger, port int) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})

	swaggerUrl := fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(httpSwagger.URL(swaggerUrl)))
	handler := v1.NewHandler(services, logger)
	handler.InitAPI(mux)
	return mux
}
