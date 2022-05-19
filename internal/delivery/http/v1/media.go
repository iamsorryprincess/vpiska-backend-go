package v1

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

func (h *Handler) initMediaAPI(router *gin.RouterGroup) {
	media := router.Group("/media")
	media.POST("", h.uploadMedia)
	media.GET("/:id", h.getMedia)
	media.DELETE("/:id", h.deleteMedia)
	media.GET("/metadata/:id", h.getMetadata)
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
func (h *Handler) uploadMedia(context *gin.Context) {
	fileData, header, err := parseFormFile("file", context, h.errorLogger)

	if err != nil {
		return
	}

	mediaId, err := h.services.Media.Create(context.Request.Context(), &service.CreateMediaInput{
		Name:        header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		Data:        fileData,
	})

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	writeResponse(mediaId, context)
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
func (h *Handler) getMedia(context *gin.Context) {
	mediaId := context.Param("id")

	if mediaId == "" {
		context.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	isMatched, err := regexp.MatchString(idRegexp, mediaId)

	if err != nil {
		h.errorLogger.Println(err)
		context.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isMatched {
		context.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaData, err := h.services.Media.GetFile(context.Request.Context(), mediaId)

	if err != nil {
		if err == domain.ErrMediaNotFound {
			context.Writer.WriteHeader(http.StatusNotFound)
			return
		}

		h.errorLogger.Println(err)
		context.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	context.Writer.Header().Set("Content-Type", mediaData.ContentType)
	context.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", mediaData.Size))
	context.Writer.WriteHeader(http.StatusOK)
	_, err = context.Writer.Write(mediaData.Data)

	if err != nil {
		h.errorLogger.Println(err)
		context.Writer.WriteHeader(http.StatusInternalServerError)
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
func (h *Handler) getMetadata(context *gin.Context) {
	mediaId := context.Param("id")

	if mediaId == "" {
		response := createDomainErrorResponse(errEmptyId)
		context.JSON(http.StatusOK, response)
		return
	}

	isMatched, err := regexp.MatchString(idRegexp, mediaId)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	if !isMatched {
		response := createDomainErrorResponse(errInvalidId)
		context.JSON(http.StatusOK, response)
		return
	}

	metadata, err := h.services.Media.GetMetadata(context.Request.Context(), mediaId)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	writeResponse(fileMetadataResponse{
		ID:          metadata.ID,
		Name:        metadata.Name,
		Size:        metadata.Size,
		ContentType: metadata.ContentType,
	}, context)
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
func (h *Handler) deleteMedia(context *gin.Context) {
	mediaId := context.Param("id")

	if mediaId == "" {
		writeErrorResponse(errEmptyId, h.errorLogger, context)
		return
	}

	isMatched, err := regexp.MatchString(idRegexp, mediaId)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	if !isMatched {
		response := createDomainErrorResponse(errInvalidId)
		context.JSON(http.StatusOK, response)
		return
	}

	err = h.services.Media.Delete(context.Request.Context(), mediaId)

	if err != nil {
		writeErrorResponse(err, h.errorLogger, context)
		return
	}

	writeResponse(nil, context)
}
