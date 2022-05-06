package v1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const contentTypeJSON = "application/json"

var (
	invalidMethodError      = errors.New("invalid method")
	invalidContentTypeError = errors.New("invalid Content-Type")
	validationError         = errors.New("invalid request")
)

type errorResponse struct {
	ErrorCode string `json:"errorCode"`
}

type apiResponse struct {
	IsSuccess bool             `json:"isSuccess"`
	Result    interface{}      `json:"result"`
	Errors    []*errorResponse `json:"errors"`
}

type validated interface {
	Validate() ([]string, error)
}

func createDomainErrorResponse(err error) *apiResponse {
	return &apiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors: []*errorResponse{{
			ErrorCode: err.Error(),
		}},
	}
}

func createValidationErrorResponse(errs []string) *apiResponse {
	errorsResponses := make([]*errorResponse, len(errs))

	for index, item := range errs {
		errorsResponses[index] = &errorResponse{
			ErrorCode: item,
		}
	}

	return &apiResponse{
		IsSuccess: false,
		Result:    nil,
		Errors:    errorsResponses,
	}
}

func createSuccessResponse(response interface{}) *apiResponse {
	return &apiResponse{
		IsSuccess: true,
		Result:    response,
		Errors:    nil,
	}
}

func parseRequest(writer http.ResponseWriter, request *http.Request, method string, contentType string) ([]byte, error) {
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

func validateRequest(writer http.ResponseWriter, request validated) error {
	validationErrors, err := request.Validate()

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return err
	}

	if validationErrors != nil {
		errResponse := createValidationErrorResponse(validationErrors)
		bytes, marshalErr := json.Marshal(errResponse)

		if marshalErr != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return marshalErr
		}

		writer.Header().Set("Content-Type", contentTypeJSON)
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
	errResponse := createDomainErrorResponse(domainError)
	bytes, marshalErr := json.Marshal(errResponse)

	if marshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return marshalErr
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
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

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, writeErr := writer.Write(bytes)

	if writeErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return writeErr
	}

	return nil
}

func deserializeAndValidateRequest(writer http.ResponseWriter, data []byte, reqBody validated) error {
	if unMarshalErr := json.Unmarshal(data, reqBody); unMarshalErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return unMarshalErr
	}

	return validateRequest(writer, reqBody)
}

func writeResponse(writer http.ResponseWriter, response interface{}, domainError error) {
	if domainError != nil {
		_ = writeDomainErrorResponse(writer, domainError)
		return
	}

	_ = writeSuccessResponse(writer, response)
}
