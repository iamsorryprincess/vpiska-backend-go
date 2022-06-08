package v1

import (
	"net/http"
	"regexp"

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

func (h *Handler) initUsersAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/users/create", h.createUser)
	mux.HandleFunc("/api/v1/users/login", h.loginUser)
	mux.HandleFunc("/api/v1/users/password/change", h.jwtAuth(h.changePassword))
	mux.HandleFunc("/api/v1/users/update", h.jwtAuth(h.updateUser))
	mux.HandleFunc("/api/v1/users/media/set", h.jwtAuth(h.setUserImage))
}

type loginResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	ImageID     *string `json:"imageId"`
	EventID     *string `json:"eventId"`
	AccessToken string  `json:"accessToken"`
}

type createUserRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (r createUserRequest) Validate() ([]string, error) {
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
	reqBody := createUserRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	response, err := h.services.Users.Create(request.Context(), service.CreateUserInput{
		Name:     reqBody.Name,
		Phone:    reqBody.Phone,
		Password: reqBody.Password,
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(response))
}

type loginUserRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (r loginUserRequest) Validate() ([]string, error) {
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
	reqBody := loginUserRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	response, err := h.services.Users.Login(request.Context(), service.LoginUserInput{
		Phone:    reqBody.Phone,
		Password: reqBody.Password,
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(response))
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
}

type changePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (r changePasswordRequest) Validate() ([]string, error) {
	var validationErrors []string

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

// ChangePassword godoc
// @Summary      Изменить пароль
// @Security     UserAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body changePasswordRequest true "body"
// @Success      200 {object} apiResponse{result=tokenResponse}
// @Router       /v1/users/password/change [post]
func (h *Handler) changePassword(writer http.ResponseWriter, request *http.Request) {
	reqBody := changePasswordRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	result, err := h.services.Users.ChangePassword(request.Context(), service.ChangePasswordInput{
		ID:       getUserID(request),
		Password: reqBody.Password,
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(tokenResponse{
		AccessToken: result,
	}))
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (r updateUserRequest) Validate() ([]string, error) {
	var validationErrors []string

	if r.Phone != "" {
		if matched, err := regexp.MatchString(phoneRegexp, r.Phone); err != nil {
			return nil, err
		} else if !matched {
			validationErrors = append(validationErrors, invalidPhoneFormatError)
		}
	}

	return validationErrors, nil
}

// UpdateUser godoc
// @Summary      Обновить информацию о пользователе
// @Security     UserAuth
// @Tags         users
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body updateUserRequest false "body"
// @Success      200 {object} apiResponse{result=tokenResponse}
// @Router       /v1/users/update [post]
func (h *Handler) updateUser(writer http.ResponseWriter, request *http.Request) {
	reqBody := updateUserRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	result, err := h.services.Users.Update(request.Context(), service.UpdateUserInput{
		ID:    getUserID(request),
		Name:  reqBody.Name,
		Phone: reqBody.Phone,
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(tokenResponse{
		AccessToken: result,
	}))
}

type setImageResponse struct {
	ImageID     string `json:"imageId"`
	AccessToken string `json:"accessToken"`
}

// SetUserImage godoc
// @Summary      Установить пользователю картинку
// @Security     UserAuth
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Content-Type application/json
// @param        image formData file true "file"
// @Success      200 {object} apiResponse{result=setImageResponse}
// @Router       /v1/users/media/set [post]
func (h *Handler) setUserImage(writer http.ResponseWriter, request *http.Request) {
	file, header, isOk := h.parseForm(writer, request, "image")

	if !isOk {
		return
	}

	data := make([]byte, header.Size)

	if _, err := file.Read(data); err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return
	}

	imageId, accessToken, err := h.services.Users.SetUserImage(request.Context(), &service.SetUserImageInput{
		UserID:      getUserID(request),
		FileName:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		FileData:    data,
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(setImageResponse{
		ImageID:     imageId,
		AccessToken: accessToken,
	}))
}
