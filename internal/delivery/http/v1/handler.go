package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const (
	contentTypeJSON   = "application/json"
	multipartFormData = "multipart/form-data"
)

var (
	errInvalidMethod      = errors.New("invalid method")
	errInvalidContentType = errors.New("invalid Content-Type")
	errInvalidRequest     = errors.New("invalid request")
)

type Handler struct {
	services *service.Services
	logger   *log.Logger
}

func NewHandler(services *service.Services, logger *log.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitAPI(mux *http.ServeMux) http.Handler {
	recoveryMiddleware := newRecoveryMiddleware(h.logger)
	h.initUsersAPI(mux)
	h.initMediaAPI(mux)
	return recoveryMiddleware(mux)
}

func (h *Handler) writeErrorResponse(writer http.ResponseWriter, handleErr error) {
	response := createDomainErrorResponse(handleErr)
	bytes, err := json.Marshal(&response)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
		return
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(bytes)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
	}
}

func (h *Handler) handlePostJSON(writer http.ResponseWriter, request *http.Request) ([]byte, error) {
	if request.Method != http.MethodPost {
		h.writeErrorResponse(writer, errInvalidMethod)
		return nil, errInvalidMethod
	}

	if request.Header.Get("Content-Type") != contentTypeJSON {
		h.writeErrorResponse(writer, errInvalidContentType)
		return nil, errInvalidContentType
	}

	data, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, domain.ErrInternalError)
		return nil, err
	}

	return data, nil
}

func (h *Handler) bindJSON(data []byte, request requestForValidate, writer http.ResponseWriter) error {
	if err := json.Unmarshal(data, request); err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, domain.ErrInternalError)
		return err
	}

	validationErrs, err := request.validate()

	if err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, domain.ErrInternalError)
		return err
	}

	if validationErrs != nil {
		response := createValidationErrorResponse(validationErrs)
		bytes, err := json.Marshal(&response)

		if err != nil {
			h.logger.Println(err)
			h.writeErrorResponse(writer, domain.ErrInternalError)
			return err
		}

		writer.Header().Set("Content-Type", contentTypeJSON)
		writer.WriteHeader(http.StatusOK)
		_, err = writer.Write(bytes)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			h.logger.Println(err)
			return err
		}

		return errInvalidRequest
	}

	return nil
}

func (h *Handler) getFormFile(writer http.ResponseWriter, request *http.Request, filename string) (multipart.File, *multipart.FileHeader, error) {
	if !strings.Contains(request.Header.Get("Content-Type"), multipartFormData) {
		h.writeErrorResponse(writer, errInvalidContentType)
		return nil, nil, errInvalidContentType
	}

	err := request.ParseMultipartForm(10 << 20)

	if err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, domain.ErrInternalError)
		return nil, nil, err
	}

	if request.MultipartForm.File[filename] == nil {
		h.writeErrorResponse(writer, domain.ErrEmptyMedia)
		return nil, nil, domain.ErrEmptyMedia
	}

	file, header, err := request.FormFile(filename)

	if err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, domain.ErrInternalError)
	}

	return file, header, err
}

func (h *Handler) writeDomainErrorResponse(writer http.ResponseWriter, domainError error) {
	mappedErr := domain.MapDomainError(domainError)
	response := createDomainErrorResponse(mappedErr)
	bytes, err := json.Marshal(&response)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
		return
	}

	if mappedErr == domain.ErrInternalError {
		h.logger.Println(domainError)
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(bytes)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
	}
}

func (h *Handler) writeSuccessResponse(writer http.ResponseWriter, response interface{}) {
	result := createSuccessResponse(response)
	bytes, err := json.Marshal(&result)

	if err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, err)
		return
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(bytes)

	if err != nil {
		h.logger.Println(err)
		h.writeErrorResponse(writer, err)
	}
}
