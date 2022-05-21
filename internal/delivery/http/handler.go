package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(services *service.Services, logger logger.Logger, tokenManager auth.TokenManager, port int) http.Handler {
	gin.SetMode("release")
	ginEngine := gin.New()

	ginEngine.GET("/health", func(context *gin.Context) {
		context.Writer.WriteHeader(http.StatusOK)
	})

	swaggerUrl := fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)
	swaggerHandler := httpSwagger.Handler(httpSwagger.URL(swaggerUrl))
	ginEngine.GET("/swagger/*any", gin.WrapH(swaggerHandler))
	apiRouter := ginEngine.Group("/api")
	handler := v1.NewHandler(logger, services, tokenManager)
	handler.InitAPI(apiRouter)
	return ginEngine
}
