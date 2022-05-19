package v1

import (
	"errors"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
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
	errorLogger  *log.Logger
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(errorLogger *log.Logger, services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		errorLogger:  errorLogger,
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) InitAPI(router *gin.RouterGroup) {
	v1Router := router.Group("/v1")
	h.initUsersAPI(v1Router)
	h.initMediaAPI(v1Router)
}

func parseFormFile(name string, context *gin.Context, logger *log.Logger) ([]byte, *multipart.FileHeader, error) {
	header, err := context.FormFile(name)

	if err != nil {
		if err.Error() == emptyFormFileError {
			response := createDomainErrorResponse(errEmptyMedia)
			context.JSON(http.StatusOK, response)
			return nil, nil, err
		}

		writeErrorResponse(err, logger, context)
		return nil, nil, err
	}

	file, err := header.Open()

	if err != nil {
		writeErrorResponse(err, logger, context)
		return nil, nil, err
	}

	defer file.Close()
	fileData := make([]byte, header.Size)
	_, err = file.Read(fileData)

	if err != nil {
		writeErrorResponse(err, logger, context)
		return nil, nil, err
	}

	return fileData, header, nil
}
