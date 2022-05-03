package v1

import (
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

type CreateUserRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func CreateUserHandler(userService *service.UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, parseError := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

		if parseError != nil {
			return
		}

		reqBody := &CreateUserRequest{}

		if validError := deserializeAndValidateRequest(writer, data, reqBody); validError != nil {
			return
		}

		response, domainError := userService.Create(request.Context(), reqBody.Name, reqBody.Phone, reqBody.Password)
		writeResponse(writer, response, domainError)
	}
}

type LoginUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func LoginUserHandler(userService *service.UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, parseError := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

		if parseError != nil {
			return
		}

		reqBody := &LoginUserRequest{}

		if validError := deserializeAndValidateRequest(writer, data, reqBody); validError != nil {
			return
		}

		response, domainError := userService.Login(request.Context(), reqBody.Phone, reqBody.Password)
		writeResponse(writer, response, domainError)
	}
}

type ChangePasswordRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func ChangePasswordHandler(userService *service.UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		data, parseError := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

		if parseError != nil {
			return
		}

		reqBody := &ChangePasswordRequest{}

		if validError := deserializeAndValidateRequest(writer, data, reqBody); validError != nil {
			return
		}

		response, domainError := userService.ChangePassword(request.Context(), reqBody.ID, reqBody.Password)
		writeResponse(writer, response, domainError)
	}
}
