package v1

import (
	"net/http"
	"regexp"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func (h *Handler) initUsersApi(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/users/create", h.createUser)
	mux.HandleFunc("/api/v1/users/login", h.loginUser)
	mux.HandleFunc("/api/v1/users/password/change", h.changePassword)
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
// @Router       /v1/users/create [post]
func (h *Handler) createUser(writer http.ResponseWriter, request *http.Request) {
	data, handleErr := h.handlePostJSON(writer, request)

	if handleErr != nil {
		return
	}

	body := createUserRequest{}
	bindErr := h.bindJSON(data, &body, writer)

	if bindErr != nil {
		return
	}

	input := body.toServiceInput()
	result, domainErr := h.services.Users.Create(request.Context(), input)

	if domainErr != nil {
		h.writeDomainErrorResponse(writer, domainErr)
		return
	}

	response := mapToResponse(result)
	h.writeSuccessResponse(writer, response)
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
// @Router       /v1/users/login [post]
func (h *Handler) loginUser(writer http.ResponseWriter, request *http.Request) {
	data, handleErr := h.handlePostJSON(writer, request)

	if handleErr != nil {
		return
	}

	body := loginUserRequest{}
	bindErr := h.bindJSON(data, &body, writer)

	if bindErr != nil {
		return
	}

	input := body.toServiceInput()
	result, domainErr := h.services.Users.Login(request.Context(), input)

	if domainErr != nil {
		h.writeDomainErrorResponse(writer, domainErr)
		return
	}

	response := mapToResponse(result)
	h.writeSuccessResponse(writer, response)
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
// @Router       /v1/users/password/change [post]
func (h *Handler) changePassword(writer http.ResponseWriter, request *http.Request) {
	data, handleErr := h.handlePostJSON(writer, request)

	if handleErr != nil {
		return
	}

	body := changePasswordRequest{}
	bindErr := h.bindJSON(data, &body, writer)

	if bindErr != nil {
		return
	}

	input := body.toServiceInput()
	result, domainErr := h.services.Users.ChangePassword(request.Context(), input)

	if domainErr != nil {
		h.writeDomainErrorResponse(writer, domainErr)
		return
	}

	response := mapToResponse(result)
	h.writeSuccessResponse(writer, response)
}

func (r *createUserRequest) validate() ([]string, error) {
	var validationErrors []string

	if r.Name == "" {
		validationErrors = append(validationErrors, emptyNameError)
	}

	if r.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, r.Phone); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if r.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(r.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	if r.Password != r.ConfirmPassword {
		validationErrors = append(validationErrors, invalidConfirmPasswordError)
	}

	return validationErrors, nil
}

func (r *createUserRequest) toServiceInput() service.CreateUserInput {
	return service.CreateUserInput{
		Name:     r.Name,
		Phone:    r.Phone,
		Password: r.Password,
	}
}

func (r *loginUserRequest) validate() ([]string, error) {
	var validationErrors []string

	if r.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, r.Phone); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if r.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(r.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	return validationErrors, nil
}

func (r *loginUserRequest) toServiceInput() service.LoginUserInput {
	return service.LoginUserInput{
		Phone:    r.Phone,
		Password: r.Password,
	}
}

func (r *changePasswordRequest) validate() ([]string, error) {
	var validationErrors []string

	if r.ID == "" {
		validationErrors = append(validationErrors, emptyIDError)
	} else if matched, err := regexp.MatchString(idRegexp, r.ID); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidIdFormatError)
	}

	if r.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(r.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	if r.Password != r.ConfirmPassword {
		validationErrors = append(validationErrors, invalidConfirmPasswordError)
	}

	return validationErrors, nil
}

func (r *changePasswordRequest) toServiceInput() service.ChangePasswordInput {
	return service.ChangePasswordInput{
		ID:       r.ID,
		Password: r.Password,
	}
}

func mapToResponse(response service.LoginResponse) loginResponse {
	return loginResponse{
		ID:          response.ID,
		Name:        response.Name,
		Phone:       response.Phone,
		ImageID:     response.ImageID,
		AccessToken: response.AccessToken,
	}
}
