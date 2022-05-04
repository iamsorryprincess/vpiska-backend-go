package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const ContentTypeJSON = "application/json"

var (
	invalidMethodError      = errors.New("invalid method")
	invalidContentTypeError = errors.New("invalid Content-Type")
	validationError         = errors.New("invalid request")
)

type Validated interface {
	Validate() ([]string, error)
}

func ParseRequest(writer http.ResponseWriter, request *http.Request, method string, contentType string) ([]byte, error) {
	if request.Method != method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return nil, invalidMethodError
	}

	if request.Header.Get("Content-Type") != contentType {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return nil, invalidContentTypeError
	}

	data, readErr := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if readErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return nil, readErr
	}

	return data, nil
}

func DeserializeAndValidateRequest(writer http.ResponseWriter, data []byte, reqBody Validated) error {
	if unMarshalErr := json.Unmarshal(data, reqBody); unMarshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return unMarshalErr
	}

	return validateRequest(writer, reqBody)
}

func WriteResponse(writer http.ResponseWriter, response interface{}, domainError error) {
	if domainError != nil {
		_ = writeDomainErrorResponse(writer, domainError)
		return
	}

	_ = writeSuccessResponse(writer, response)
}

func validateRequest(writer http.ResponseWriter, request Validated) error {
	validationErrors, err := request.Validate()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return err
	}

	if validationErrors != nil {
		errorResponse := createValidationErrorResponse(validationErrors)
		bytes, marshalErr := json.Marshal(errorResponse)

		if marshalErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return marshalErr
		}

		writer.Header().Set("Content-Type", ContentTypeJSON)
		writer.WriteHeader(http.StatusOK)
		_, writeErr := writer.Write(bytes)

		if writeErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return writeErr
		}

		return validationError
	}

	return nil
}

func writeDomainErrorResponse(writer http.ResponseWriter, domainError error) error {
	errorResponse := createDomainErrorResponse(domainError)
	bytes, marshalErr := json.Marshal(errorResponse)

	if marshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return marshalErr
	}

	writer.Header().Set("Content-Type", ContentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return writeErr
	}

	return nil
}

func writeSuccessResponse(writer http.ResponseWriter, response interface{}) error {
	bytes, marshalErr := json.Marshal(createSuccessResponse(response))

	if marshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return marshalErr
	}

	writer.Header().Set("Content-Type", ContentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return writeErr
	}

	return nil
}
