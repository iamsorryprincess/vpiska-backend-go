package v1

import (
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

type Handler struct {
	logger       logger.Logger
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(logger logger.Logger, services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		logger:       logger,
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitAPI(mux *http.ServeMux) {
	h.initUsersAPI(mux)
}
