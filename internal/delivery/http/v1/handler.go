package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/logger"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

type Handler struct {
	services *service.Services
	logger   logger.Logger
}

func NewHandler(services *service.Services, logger logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitAPI(mux *http.ServeMux) {
	h.initUsersApi(mux)
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

	data, readErr := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if readErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(readErr)
		return nil, readErr
	}

	return data, nil
}

func (h *Handler) bindJSON(data []byte, request requestForValidate, writer http.ResponseWriter) error {
	if unMarshalErr := json.Unmarshal(data, request); unMarshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(unMarshalErr)
		return unMarshalErr
	}

	validationErrs, err := request.validate()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(err)
		return err
	}

	if validationErrs != nil {
		errResponse := createValidationErrorResponse(validationErrs)
		bytes, marshalErr := json.Marshal(errResponse)

		if marshalErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			h.logger.LogError(marshalErr)
			return marshalErr
		}

		writer.Header().Set("Content-Type", contentTypeJSON)
		writer.WriteHeader(http.StatusOK)
		_, writeErr := writer.Write(bytes)

		if writeErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			h.logger.LogError(writeErr)
			return writeErr
		}

		return errInvalidRequest
	}

	return nil
}

func (h *Handler) writeDomainErrorResponse(writer http.ResponseWriter, domainError error) {
	errResponse := createDomainErrorResponse(domainError)
	bytes, marshalErr := json.Marshal(errResponse)

	if marshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(marshalErr)
		return
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(marshalErr)
		return
	}
}

func (h *Handler) writeSuccessResponse(writer http.ResponseWriter, response interface{}) {
	bytes, marshalErr := json.Marshal(createSuccessResponse(response))

	if marshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(marshalErr)
		return
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.LogError(marshalErr)
		return
	}
}
