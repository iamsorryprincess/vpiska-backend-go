package v1

import (
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

type CreateUserRequest struct {
	Name            string
	Phone           string
	Password        string
	ConfirmPassword string
}

func CreateUserHandler(userService *service.UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, isParsed := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

		if !isParsed {
			return
		}

		reqBody := &CreateUserRequest{}

		if valid := deserializeAndValidateRequest(writer, data, reqBody); !valid {
			return
		}

		response, domainError := userService.Create(request.Context(), reqBody.Name, reqBody.Phone, reqBody.Password)
		writeResponse(writer, response, domainError)
	}
}

type LoginUserRequest struct {
	Phone    string
	Password string
}

func LoginUserHandler(userService *service.UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, isParsed := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

		if !isParsed {
			return
		}

		reqBody := &LoginUserRequest{}

		if valid := deserializeAndValidateRequest(writer, data, reqBody); !valid {
			return
		}

		response, domainError := userService.Login(request.Context(), reqBody.Phone, reqBody.Password)
		writeResponse(writer, response, domainError)
	}
}
