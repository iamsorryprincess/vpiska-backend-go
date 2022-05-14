package v1

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const (
	idRegexp             = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`
	emptyIDError         = "IdIsEmpty"
	invalidIdFormatError = "InvalidIdFormat"
	emptyFormFileError   = "http: no such file"
)

var (
	errEmptyId    = errors.New(emptyIDError)
	errInvalidId  = errors.New(invalidIdFormatError)
	errEmptyMedia = errors.New("param file is empty")
)

type Handler struct {
	errorLogger *log.Logger
	services    *service.Services
}

func NewHandler(errorLogger *log.Logger, services *service.Services) *Handler {
	return &Handler{
		errorLogger: errorLogger,
		services:    services,
	}
}

func (h *Handler) InitAPI(router *gin.RouterGroup) {
	v1Router := router.Group("/v1")
	h.initUsersAPI(v1Router)
	h.initMediaAPI(v1Router)
}
