package handler

import (
	"net/http"

	v1 "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1"
	req "github.com/iamsorryprincess/vpiska-backend-go/internal/delivery/http/v1/request"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

func (h *UserHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	data, parseError := v1.ParseRequest(writer, request, http.MethodPost, v1.ContentTypeJSON)

	if parseError != nil {
		return
	}

	reqBody := &req.CreateUserRequest{}

	if validError := v1.DeserializeAndValidateRequest(writer, data, reqBody); validError != nil {
		return
	}

	response, domainError := h.service.Create(request.Context(), reqBody.Name, reqBody.Phone, reqBody.Password)
	v1.WriteResponse(writer, response, domainError)
}

func (h *UserHandler) LoginUser(writer http.ResponseWriter, request *http.Request) {
	data, parseError := v1.ParseRequest(writer, request, http.MethodPost, v1.ContentTypeJSON)

	if parseError != nil {
		return
	}

	reqBody := &req.LoginUserRequest{}

	if validError := v1.DeserializeAndValidateRequest(writer, data, reqBody); validError != nil {
		return
	}

	response, domainError := h.service.Login(request.Context(), reqBody.Phone, reqBody.Password)
	v1.WriteResponse(writer, response, domainError)
}

func (h *UserHandler) ChangePassword(writer http.ResponseWriter, request *http.Request) {
	data, parseError := v1.ParseRequest(writer, request, http.MethodPost, v1.ContentTypeJSON)

	if parseError != nil {
		return
	}

	reqBody := &req.ChangePasswordRequest{}

	if validError := v1.DeserializeAndValidateRequest(writer, data, reqBody); validError != nil {
		return
	}

	response, domainError := h.service.ChangePassword(request.Context(), reqBody.ID, reqBody.Password)
	v1.WriteResponse(writer, response, domainError)
}
