package v1

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
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

func (h *Handler) handlePostJSON(writer http.ResponseWriter, request *http.Request) ([]byte, error) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return nil, errInvalidMethod
	}

	if request.Header.Get("Content-Type") != contentTypeJSON {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return nil, errInvalidContentType
	}

	data, err := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
		return nil, err
	}

	return data, nil
}

func (h *Handler) bindJSON(data []byte, request requestForValidate, writer http.ResponseWriter) error {
	if err := json.Unmarshal(data, request); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
		return err
	}

	validationErrs, err := request.validate()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
		return err
	}

	if validationErrs != nil {
		response := createValidationErrorResponse(validationErrs)
		bytes, err := json.Marshal(&response)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			h.logger.Println(err)
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
