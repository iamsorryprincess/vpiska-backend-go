package v1

import (
	"net/http"

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

type loginResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	ImageID     string `json:"imageId"`
	AccessToken string `json:"accessToken"`
}

type createUserRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// CreateUser godoc
// @Summary      Создать пользователя
// @Tags         users
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body createUserRequest true "body"
// @Success      200 {object} apiResponse{result=loginResponse}
// @Router       /users/create [post]
func (h *UserHandler) CreateUser(writer http.ResponseWriter, request *http.Request) {
	data, parseError := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

	if parseError != nil {
		return
	}

	reqBody := &createUserRequest{}

	if validError := deserializeAndValidateRequest(writer, data, reqBody); validError != nil {
		return
	}

	response, domainError := h.service.Create(request.Context(), reqBody.Name, reqBody.Phone, reqBody.Password)
	writeResponse(writer, response, domainError)
}

type loginUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// LoginUser godoc
// @Summary      Войти в систему
// @Tags         users
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body loginUserRequest true "body"
// @Success      200 {object} apiResponse{result=loginResponse}
// @Router       /users/login [post]
func (h *UserHandler) LoginUser(writer http.ResponseWriter, request *http.Request) {
	data, parseError := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

	if parseError != nil {
		return
	}

	reqBody := &loginUserRequest{}

	if validError := deserializeAndValidateRequest(writer, data, reqBody); validError != nil {
		return
	}

	response, domainError := h.service.Login(request.Context(), reqBody.Phone, reqBody.Password)
	writeResponse(writer, response, domainError)
}

type changePasswordRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ChangePassword godoc
// @Summary      Изменить пароль
// @Tags         users
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body changePasswordRequest true "body"
// @Success      200 {object} apiResponse{result=loginResponse}
// @Router       /users/password/change [post]
func (h *UserHandler) ChangePassword(writer http.ResponseWriter, request *http.Request) {
	data, parseError := parseRequest(writer, request, http.MethodPost, contentTypeJSON)

	if parseError != nil {
		return
	}

	reqBody := &changePasswordRequest{}

	if validError := deserializeAndValidateRequest(writer, data, reqBody); validError != nil {
		return
	}

	response, domainError := h.service.ChangePassword(request.Context(), reqBody.ID, reqBody.Password)
	writeResponse(writer, response, domainError)
}
