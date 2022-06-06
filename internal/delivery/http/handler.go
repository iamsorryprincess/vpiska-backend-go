package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(services *service.Services, logger logger.Logger, tokenManager auth.TokenManager) http.Handler {
	gin.SetMode("release")
	ginEngine := gin.New()

	ginEngine.GET("/health", func(context *gin.Context) {
		context.Writer.WriteHeader(http.StatusOK)
	})

	swaggerHandler := httpSwagger.Handler(httpSwagger.URL("doc.json"))
	ginEngine.GET("/swagger/*any", gin.WrapH(swaggerHandler))
	apiRouter := ginEngine.Group("/api")
	handler := v1.NewHandler(logger, services, tokenManager)
	handler.InitAPI(apiRouter)
	websocketsHandler := websocket.NewHandler(logger, tokenManager, services.Events, services.Publisher)
	ginEngine.GET("/api/v1/websockets/*any", gin.WrapH(websocketsHandler))
	return ginEngine
}
