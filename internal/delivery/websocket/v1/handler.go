package v1

import (
	"net/http"
	"time"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

type Handler struct {
	pingPeriod   time.Duration
	logger       logger.Logger
	tokenManager auth.TokenManager
	events       service.Events
	publisher    service.Publisher
}

func NewHandler(pingPeriod time.Duration,
	logger logger.Logger,
	tokenManager auth.TokenManager,
	events service.Events,
	publisher service.Publisher) *Handler {
	return &Handler{
		pingPeriod:   pingPeriod,
		logger:       logger,
		tokenManager: tokenManager,
		events:       events,
		publisher:    publisher,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/websockets/event", h.upgradeEventConnection)
}
