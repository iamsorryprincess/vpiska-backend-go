package websocket

import (
	"net/http"

	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/websocket/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

func NewHandler(logger logger.Logger, tokenManager auth.TokenManager, events service.Events, publisher service.Publisher) http.Handler {
	mux := http.NewServeMux()
	v1Handler := v1.NewHandler(logger, tokenManager, events, publisher)
	v1Handler.InitRoutes(mux)
	return mux
}
