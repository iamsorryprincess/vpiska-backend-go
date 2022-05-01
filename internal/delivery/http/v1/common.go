package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const contentTypeJSON = "application/json"

func parseRequest(writer http.ResponseWriter, request *http.Request, method string, contentType string) ([]byte, bool) {
	if request.Method != method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return nil, false
	}

	if request.Header.Get("Content-Type") != contentType {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return nil, false
	}

	data, readErr := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if readErr != nil {
		panic(readErr)
	}

	return data, true
}

func deserializeAndValidateRequest(writer http.ResponseWriter, data []byte, reqBody Validated) bool {
	if unMarshalErr := json.Unmarshal(data, reqBody); unMarshalErr != nil {
		panic(unMarshalErr)
	}

	return validateRequest(writer, reqBody)
}

func writeResponse(writer http.ResponseWriter, response interface{}, domainError error) {
	if domainError != nil {
		writeDomainErrorResponse(writer, domainError)
		return
	}

	writeSuccessResponse(writer, response)
}
