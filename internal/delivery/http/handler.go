package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHandler(services *service.Services, logger *log.Logger, tokenManager auth.TokenManager, port int) http.Handler {
	ginEngine := gin.Default()

	ginEngine.Use(gin.Logger())
	ginEngine.Use(gin.Recovery())

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
