package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const (
	emptyAddressError         = "AddressIsEmpty"
	emptyCoordinatesError     = "CoordinatesIsEmpty"
	emptyHorizontalRangeError = "HorizontalRangeIsEmpty"
	emptyVerticalRangeError   = "VerticalRangeIsEmpty"
)

func (h *Handler) initEventsAPI(router *gin.RouterGroup) {
	events := router.Group("/events")
	events.POST("/get", h.getEventByID)
	events.POST("/range", h.getEventsByRange)
	authenticated := events.Group("/", h.jwtAuth)
	authenticated.POST("/create", h.createEvent)
	authenticated.POST("/close", h.closeEvent)
}

type coordinates struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
}

type mediaInfo struct {
	ID          string `json:"id"`
	ContentType string `json:"contentType"`
}

type chatMessage struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	UserImageID string `json:"userImageId"`
	Message     string `json:"message"`
}

type eventResponse struct {
	ID           string        `json:"id"`
	OwnerID      string        `json:"ownerId"`
	Name         string        `json:"name"`
	Address      string        `json:"address"`
	Coordinates  coordinates   `json:"coordinates"`
	UsersCount   int           `json:"usersCount"`
	Media        []mediaInfo   `json:"media"`
	ChatMessages []chatMessage `json:"chatMessages"`
}

type createEventRequest struct {
	Name        string       `json:"name"`
	Address     string       `json:"address"`
	Coordinates *coordinates `json:"coordinates"`
}

// CreateEvent godoc
// @Summary      Создать эвент
// @Security     UserAuth
// @Tags         events
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body createEventRequest true "body"
// @Success      200 {object} apiResponse{result=eventResponse}
// @Router       /v1/events/create [post]
func (h *Handler) createEvent(context *gin.Context) {
	request := createEventRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	validationErrors := validateCreateEventRequest(request)

	if len(validationErrors) > 0 {
		writeValidationErrResponse(validationErrors, context)
		return
	}

	result, err := h.services.Events.Create(context.Request.Context(), service.CreateEventInput{
		OwnerID: context.GetString("UserID"),
		Name:    request.Name,
		Address: request.Address,
		Coordinates: domain.Coordinates{
			X: *request.Coordinates.X,
			Y: *request.Coordinates.Y,
		},
	})

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	writeResponse(result, context)
}

type eventIDRequest struct {
	EventID string `json:"eventId"`
}

// GetEvent godoc
// @Summary      Получить эвент по идентификатору
// @Tags         events
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body eventIDRequest true "body"
// @Success      200 {object} apiResponse{result=eventResponse}
// @Router       /v1/events/get [post]
func (h *Handler) getEventByID(context *gin.Context) {
	request := eventIDRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	validationErr, err := validateId(request.EventID)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	if len(validationErr) > 0 {
		writeValidationErrResponse(validationErr, context)
		return
	}

	event, err := h.services.Events.GetByID(context.Request.Context(), request.EventID)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	writeResponse(event, context)
}

type getByRangeRequest struct {
	HorizontalRange *float64     `json:"horizontalRange"`
	VerticalRange   *float64     `json:"verticalRange"`
	Coordinates     *coordinates `json:"coordinates"`
}

type eventRangeData struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	UsersCount  int         `json:"usersCount"`
	Coordinates coordinates `json:"coordinates"`
}

// GetEventsByRange godoc
// @Summary      Получить эвенты по области
// @Tags         events
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body getByRangeRequest true "body"
// @Success      200 {object} apiResponse{result=eventRangeData}
// @Router       /v1/events/range [post]
func (h *Handler) getEventsByRange(context *gin.Context) {
	request := getByRangeRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	validationErrs := validateGetByRangeRequest(request)

	if len(validationErrs) > 0 {
		writeValidationErrResponse(validationErrs, context)
		return
	}

	result, err := h.services.Events.GetByRange(context.Request.Context(), service.GetByRangeInput{
		HorizontalRange: *request.HorizontalRange,
		VerticalRange:   *request.VerticalRange,
		Coordinates: domain.Coordinates{
			X: *request.Coordinates.X,
			Y: *request.Coordinates.Y,
		},
	})

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	writeResponse(result, context)
}

// CloseEvent godoc
// @Summary      Закрыть эвент
// @Security     UserAuth
// @Tags         events
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body eventIDRequest true "body"
// @Success      200 {object} apiResponse{result=string}
// @Router       /v1/events/close [post]
func (h *Handler) closeEvent(context *gin.Context) {
	request := eventIDRequest{}
	err := context.BindJSON(&request)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	validationErr, err := validateId(request.EventID)

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	if len(validationErr) > 0 {
		writeValidationErrResponse(validationErr, context)
		return
	}

	err = h.services.Events.Close(context.Request.Context(), request.EventID, context.GetString("UserID"))

	if err != nil {
		writeErrorResponse(err, h.logger, context)
		return
	}

	writeResponse(nil, context)
}

func validateCreateEventRequest(request createEventRequest) []string {
	var validationErrors []string

	if request.Name == "" {
		validationErrors = append(validationErrors, emptyNameError)
	}

	if request.Address == "" {
		validationErrors = append(validationErrors, emptyAddressError)
	}

	if request.Coordinates == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	} else if request.Coordinates.X == nil || request.Coordinates.Y == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	}

	return validationErrors
}

func validateGetByRangeRequest(request getByRangeRequest) []string {
	var validationErrors []string

	if request.HorizontalRange == nil {
		validationErrors = append(validationErrors, emptyHorizontalRangeError)
	}

	if request.VerticalRange == nil {
		validationErrors = append(validationErrors, emptyVerticalRangeError)
	}

	if request.Coordinates == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	} else if request.Coordinates.X == nil || request.Coordinates.Y == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	}

	return validationErrors
}