package websocket

import (
	"net/http"
	"time"

	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/websocket/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

func InitWebsocketsRoutes(mux *http.ServeMux, logger logger.Logger, tokenManager auth.TokenManager, events service.Events, publisher service.Publisher) {
	v1Handler := v1.NewHandler(time.Second*120, logger, tokenManager, events, publisher)
	v1Handler.InitRoutes(mux)
}
