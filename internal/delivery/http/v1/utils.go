package v1

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
)

const (
	contentTypeJSON = "application/json"
	contentTypeFORM = "multipart/form-data"
	idRegexp        = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

	internalError        = "InternalError"
	emptyIDError         = "IdIsEmpty"
	invalidIdFormatError = "InvalidIdFormat"
)

type validatedRequest interface {
	Validate() ([]string, error)
}

func (h *Handler) writeJSONResponse(writer http.ResponseWriter, response apiResponse) {
	data, err := json.Marshal(response)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", contentTypeJSON)
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(data)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) bindValidatedRequestJSON(writer http.ResponseWriter, request *http.Request, body validatedRequest) bool {
	if request.Header.Get("Content-Type") != contentTypeJSON {
		h.writeJSONResponse(writer, newErrorResponse("invalid Content-Type"))
		return false
	}

	data, err := ioutil.ReadAll(request.Body)

	if err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return false
	}

	if err = request.Body.Close(); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return false
	}

	if err = json.Unmarshal(data, body); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return false
	}

	validationErrs, err := body.Validate()

	if err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return false
	}

	if len(validationErrs) > 0 {
		h.writeJSONResponse(writer, newValidationErrsResponse(validationErrs))
		return false
	}

	return true
}

func (h *Handler) writeError(writer http.ResponseWriter, err error) {
	if domain.IsInternalError(err) {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return
	}

	h.writeJSONResponse(writer, newErrorResponse(err.Error()))
}

func (h *Handler) parseFormFile(writer http.ResponseWriter, request *http.Request, filename string) ([]byte, *multipart.FileHeader, bool) {
	if !strings.Contains(request.Header.Get("Content-Type"), contentTypeFORM) {
		h.writeJSONResponse(writer, newErrorResponse("invalid Content-Type"))
		return nil, nil, false
	}

	if err := request.ParseForm(); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return nil, nil, false
	}

	file, header, err := request.FormFile(filename)

	if err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return nil, nil, false
	}

	data := make([]byte, header.Size)

	if _, err = file.Read(data); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return nil, nil, false
	}

	if err = file.Close(); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return nil, nil, false
	}

	if err = request.Body.Close(); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return nil, nil, false
	}

	return data, header, true
}

func validateId(id string) ([]string, error) {
	var validationErrors []string

	if id == "" {
		validationErrors = append(validationErrors, emptyIDError)
	} else if matched, err := regexp.MatchString(idRegexp, id); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidIdFormatError)
	}

	return validationErrors, nil
}

func containsError(err errorResponse, errs []errorResponse) bool {
	for _, item := range errs {
		if item.ErrorCode == err.ErrorCode {
			return true
		}
	}
	return false
}
