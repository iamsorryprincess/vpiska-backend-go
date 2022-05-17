package v1

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const (
	phoneRegexp            = `^\d{10}\b$`
	requiredPasswordLength = 6
)

const (
	emptyNameError              = "NameIsEmpty"
	emptyPhoneError             = "PhoneIsEmpty"
	emptyPasswordError          = "PasswordIsEmpty"
	invalidPhoneFormatError     = "PhoneRegexInvalid"
	invalidPasswordLengthError  = "PasswordLengthInvalid"
	invalidConfirmPasswordError = "ConfirmPasswordInvalid"
)

func (h *Handler) initUsersAPI(router *gin.RouterGroup) {
	users := router.Group("/users")
	users.POST("/create", h.createUser)
	users.POST("/login", h.loginUser)
	authenticated := users.Group("/", h.jwtAuth)
	authenticated.POST("/password/change", h.changePassword)
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
func (h *Handler) createUser(context *gin.Context) {
	request := createUserRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	validationErrs, err := validateCreateRequest(request)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	if validationErrs != nil {
		writeValidationErrResponse(validationErrs, context)
		return
	}

	result, err := h.services.Users.Create(context.Request.Context(), service.CreateUserInput{
		Name:     request.Name,
		Phone:    request.Phone,
		Password: request.Password,
	})

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	writeResponse(toLoginResponse(result), context)
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
func (h *Handler) loginUser(context *gin.Context) {
	request := loginUserRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	validationErrs, err := validateLoginRequest(request)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	if validationErrs != nil {
		writeValidationErrResponse(validationErrs, context)
		return
	}

	result, err := h.services.Users.Login(context.Request.Context(), service.LoginUserInput{
		Phone:    request.Phone,
		Password: request.Password,
	})

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	writeResponse(toLoginResponse(result), context)
}

type changePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ChangePassword godoc
// @Summary      Изменить пароль
// @Security     UserAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body changePasswordRequest true "body"
// @Success      200 {object} apiResponse{result=loginResponse}
// @Router       /v1/users/password/change [post]
func (h *Handler) changePassword(context *gin.Context) {
	request := changePasswordRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	validationErrs, err := validateChangePasswordRequest(request)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	if validationErrs != nil {
		writeValidationErrResponse(validationErrs, context)
		return
	}

	result, err := h.services.Users.ChangePassword(context.Request.Context(), service.ChangePasswordInput{
		ID:       context.GetString("UserID"),
		Password: request.Password,
	})

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	writeResponse(toLoginResponse(result), context)
}

func toLoginResponse(response service.LoginResponse) loginResponse {
	return loginResponse{
		ID:          response.ID,
		Name:        response.Name,
		Phone:       response.Phone,
		ImageID:     response.ImageID,
		AccessToken: response.AccessToken,
	}
}

func validateCreateRequest(request createUserRequest) ([]string, error) {
	var validationErrors []string

	if request.Name == "" {
		validationErrors = append(validationErrors, emptyNameError)
	}

	if request.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, request.Phone); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	if request.Password != request.ConfirmPassword {
		validationErrors = append(validationErrors, invalidConfirmPasswordError)
	}

	return validationErrors, nil
}

func validateLoginRequest(request loginUserRequest) ([]string, error) {
	var validationErrors []string

	if request.Phone == "" {
		validationErrors = append(validationErrors, emptyPhoneError)
	} else if matched, err := regexp.MatchString(phoneRegexp, request.Phone); err != nil {
		return nil, err
	} else if !matched {
		validationErrors = append(validationErrors, invalidPhoneFormatError)
	}

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	return validationErrors, nil
}

func validateChangePasswordRequest(request changePasswordRequest) ([]string, error) {
	var validationErrors []string

	if request.Password == "" {
		validationErrors = append(validationErrors, emptyPasswordError)
	} else if len(request.Password) < requiredPasswordLength {
		validationErrors = append(validationErrors, invalidPasswordLengthError)
	}

	if request.Password != request.ConfirmPassword {
		validationErrors = append(validationErrors, invalidConfirmPasswordError)
	}

	return validationErrors, nil
}
