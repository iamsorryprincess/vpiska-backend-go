package v1

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func (h *Handler) initMediaAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/media", h.POST(h.uploadMedia))
	mux.HandleFunc("/api/v1/media/metadata/", h.GET(h.getMetadata))
	mux.HandleFunc("/api/v1/media/", func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			h.getMedia(writer, request)
			return
		case http.MethodDelete:
			h.deleteMedia(writer, request)
			return
		default:
			h.writeJSONResponse(writer, newErrorResponse("invalid method"))
			return
		}
	})
}

// UploadMedia godoc
// @Summary      Загрузить медиафайл
// @Tags         media
// @Accept       multipart/form-data
// @Produce      json
// @Content-Type application/json
// @param        file formData file true "file"
// @Success      200 {object} apiResponse{result=string}
// @Router       /v1/media [post]
func (h *Handler) uploadMedia(writer http.ResponseWriter, request *http.Request) {
	data, header, isOk := h.parseFormFile(writer, request, "file")

	if !isOk {
		return
	}

	mediaId, err := h.services.Media.Create(request.Context(), &service.CreateMediaInput{
		Name:        header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		Data:        data,
	})

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(mediaId))
}

// GetMedia godoc
// @Summary      Получить медиафайл
// @Tags         media
// @Accept       */*
// @Content-Type */*
// @param        id   path      string  true  "media ID"
// @Success      200
// @Failure      400
// @Failure      404
// @Router       /v1/media/{id} [get]
func (h *Handler) getMedia(writer http.ResponseWriter, request *http.Request) {
	mediaId := strings.TrimPrefix(request.RequestURI, "/api/v1/media/")

	if mediaId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	isMatched, err := regexp.MatchString(idRegexp, mediaId)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isMatched {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaData, err := h.services.Media.GetFile(request.Context(), mediaId)

	if err != nil {
		if err == domain.ErrMediaNotFound {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", mediaData.ContentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", mediaData.Size))
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(mediaData.Data)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type fileMetadataResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"contentType"`
}

// GetMetadata godoc
// @Summary      Получить метаинформацию о файле
// @Tags         media
// @Accept       */*
// @Produce      json
// @Content-Type application/json
// @param        id   path      string  true  "media ID"
// @Success      200 {object} apiResponse{result=fileMetadataResponse}
// @Router       /v1/media/metadata/{id} [get]
func (h *Handler) getMetadata(writer http.ResponseWriter, request *http.Request) {
	mediaId := strings.TrimPrefix(request.RequestURI, "/api/v1/media/metadata/")
	validationErrs, err := validateId(mediaId)

	if err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return
	}

	if len(validationErrs) > 0 {
		h.writeJSONResponse(writer, newValidationErrsResponse(validationErrs))
		return
	}

	metadata, err := h.services.Media.GetMetadata(request.Context(), mediaId)

	if err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(fileMetadataResponse{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Size:        metadata.Size,
		ContentType: metadata.ContentType,
	}))
}

// DeleteMedia godoc
// @Summary      Удалить файл
// @Tags         media
// @Accept       */*
// @Produce      json
// @Content-Type application/json
// @param        id   path      string  true  "media ID"
// @Success      200 {object} apiResponse{result=string}
// @Router       /v1/media/{id} [delete]
func (h *Handler) deleteMedia(writer http.ResponseWriter, request *http.Request) {
	mediaId := strings.TrimPrefix(request.RequestURI, "/api/v1/media/")
	validationErrs, err := validateId(mediaId)

	if err != nil {
		h.logger.LogError(err)
		h.writeJSONResponse(writer, newErrorResponse(internalError))
		return
	}

	if len(validationErrs) > 0 {
		h.writeJSONResponse(writer, newValidationErrsResponse(validationErrs))
		return
	}

	if err = h.services.Media.Delete(request.Context(), mediaId); err != nil {
		h.writeError(writer, err)
		return
	}

	h.writeJSONResponse(writer, newSuccessResponse(nil))
}
