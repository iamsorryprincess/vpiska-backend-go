package v1

import (
	"net/http"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const (
	emptyAddressError         = "AddressIsEmpty"
	emptyCoordinatesError     = "CoordinatesIsEmpty"
	emptyHorizontalRangeError = "HorizontalRangeIsEmpty"
	emptyVerticalRangeError   = "VerticalRangeIsEmpty"
)

func (h *Handler) initEventsAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/events/get", h.POST(h.getEventByID))
	mux.HandleFunc("/api/v1/events/range", h.POST(h.getEventsByRange))
	mux.HandleFunc("/api/v1/events/create", h.jwtAuth(h.POST(h.createEvent)))
	mux.HandleFunc("/api/v1/events/update", h.jwtAuth(h.POST(h.updateEvent)))
	mux.HandleFunc("/api/v1/events/close", h.jwtAuth(h.POST(h.closeEvent)))
	mux.HandleFunc("/api/v1/events/media/add", h.jwtAuth(h.POST(h.addMediaToEvent)))
	mux.HandleFunc("/api/v1/events/media/remove", h.jwtAuth(h.POST(h.removeMediaFromEvent)))
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

func (r createEventRequest) Validate() ([]string, error) {
	var validationErrors []string

	if r.Name == "" {
		validationErrors = append(validationErrors, emptyNameError)
	}

	if r.Address == "" {
		validationErrors = append(validationErrors, emptyAddressError)
	}

	if r.Coordinates == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	} else if r.Coordinates.X == nil || r.Coordinates.Y == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	}

	return validationErrors, nil
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
func (h *Handler) createEvent(writer http.ResponseWriter, request *http.Request) {
	reqBody := createEventRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	result, err := h.services.Events.Create(request.Context(), service.CreateEventInput{
		OwnerID: getUserID(request),
		Name:    reqBody.Name,
		Address: reqBody.Address,
		Coordinates: domain.Coordinates{
			X: *reqBody.Coordinates.X,
			Y: *reqBody.Coordinates.Y,
		},
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(result))
}

type eventIDRequest struct {
	EventID string `json:"eventId"`
}

func (r eventIDRequest) Validate() ([]string, error) {
	return validateId(r.EventID)
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
func (h *Handler) getEventByID(writer http.ResponseWriter, request *http.Request) {
	reqBody := eventIDRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	event, err := h.services.Events.GetByID(request.Context(), reqBody.EventID)

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(event))
}

type eventRangeData struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	UsersCount  int         `json:"usersCount"`
	Coordinates coordinates `json:"coordinates"`
}

type getByRangeRequest struct {
	HorizontalRange *float64     `json:"horizontalRange"`
	VerticalRange   *float64     `json:"verticalRange"`
	Coordinates     *coordinates `json:"coordinates"`
}

func (r getByRangeRequest) Validate() ([]string, error) {
	var validationErrors []string

	if r.HorizontalRange == nil {
		validationErrors = append(validationErrors, emptyHorizontalRangeError)
	}

	if r.VerticalRange == nil {
		validationErrors = append(validationErrors, emptyVerticalRangeError)
	}

	if r.Coordinates == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	} else if r.Coordinates.X == nil || r.Coordinates.Y == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	}

	return validationErrors, nil
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
func (h *Handler) getEventsByRange(writer http.ResponseWriter, request *http.Request) {
	reqBody := getByRangeRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	result, err := h.services.Events.GetByRange(request.Context(), service.GetByRangeInput{
		HorizontalRange: *reqBody.HorizontalRange,
		VerticalRange:   *reqBody.VerticalRange,
		Coordinates: domain.Coordinates{
			X: *reqBody.Coordinates.X,
			Y: *reqBody.Coordinates.Y,
		},
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(result))
}

type updateEventRequest struct {
	EventID     string       `json:"eventId"`
	Address     string       `json:"address"`
	Coordinates *coordinates `json:"coordinates"`
}

func (r updateEventRequest) Validate() ([]string, error) {
	validationErrors, err := validateId(r.EventID)

	if err != nil {
		return nil, err
	}

	if r.Address == "" {
		validationErrors = append(validationErrors, emptyAddressError)
	}

	if r.Coordinates == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	} else if r.Coordinates.X == nil || r.Coordinates.Y == nil {
		validationErrors = append(validationErrors, emptyCoordinatesError)
	}

	return validationErrors, nil
}

// UpdateEvent godoc
// @Summary      Обновить эвент
// @Security     UserAuth
// @Tags         events
// @Accept       json
// @Produce      json
// @Content-Type application/json
// @param        request body updateEventRequest true "body"
// @Success      200 {object} apiResponse{result=string}
// @Router       /v1/events/update [post]
func (h *Handler) updateEvent(writer http.ResponseWriter, request *http.Request) {
	reqBody := updateEventRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	if err := h.services.Events.Update(request.Context(), service.UpdateEventInput{
		UserID:  getUserID(request),
		EventID: reqBody.EventID,
		Address: reqBody.Address,
		Coordinates: domain.Coordinates{
			X: *reqBody.Coordinates.X,
			Y: *reqBody.Coordinates.Y,
		},
	}); err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(nil))
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
func (h *Handler) closeEvent(writer http.ResponseWriter, request *http.Request) {
	reqBody := eventIDRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	if err := h.services.Events.Close(request.Context(), reqBody.EventID, getUserID(request)); err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(nil))
}

// AddMediaToEvent godoc
// @Summary      добавить медиа к евенту
// @Security     UserAuth
// @Tags         events
// @Accept       multipart/form-data
// @Produce      json
// @Content-Type application/json
// @param        eventId formData string true "event id"
// @param        media formData file true "file"
// @Success      200 {object} apiResponse{result=string}
// @Router       /v1/events/media/add [post]
func (h *Handler) addMediaToEvent(writer http.ResponseWriter, request *http.Request) {
	data, header, isOk := h.parseFormFile(writer, request, "media")

	if !isOk {
		return
	}

	eventId := request.PostFormValue("eventId")
	validationErrs, err := validateId(eventId)

	if err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return
	}

	if len(validationErrs) > 0 {
		h.writeJSONResponse(writer, newValidationErrsResponse(validationErrs))
		return
	}

	if err = h.services.Events.AddMedia(request.Context(), service.AddMediaInput{
		EventID:     eventId,
		UserID:      getUserID(request),
		FileName:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		FileSize:    header.Size,
		FileData:    data,
	}); err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(nil))
}

type removeMediaRequest struct {
	EventID string `json:"eventId"`
	MediaID string `json:"mediaId"`
}

func (r removeMediaRequest) Validate() ([]string, error) {
	validationErrs, err := validateId(r.EventID)

	if err != nil {
		return nil, err
	}

	validationErrs, err = validateId(r.MediaID)
	return validationErrs, err
}

// RemoveMediaFromEvent godoc
// @Summary      удалить медиа из евента
// @Security     UserAuth
// @Tags         events
// @Accept       application/json
// @Produce      json
// @Content-Type application/json
// @param        request body removeMediaRequest true "body"
// @Success      200 {object} apiResponse{result=string}
// @Router       /v1/events/media/remove [post]
func (h *Handler) removeMediaFromEvent(writer http.ResponseWriter, request *http.Request) {
	reqBody := removeMediaRequest{}

	if isOk := h.bindValidatedRequestJSON(writer, request, &reqBody); !isOk {
		return
	}

	if err := h.services.Events.RemoveMedia(request.Context(), service.RemoveMediaInput{
		EventID: reqBody.EventID,
		MediaID: reqBody.MediaID,
		UserID:  getUserID(request),
	}); err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(nil))
}
