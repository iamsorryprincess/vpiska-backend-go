package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const getFileUrl = "/api/v1/media/"
const getFileMetadataUrl = "/api/v1/media/metadata/"

func (h *Handler) initMediaAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/media", h.uploadFile)
	mux.HandleFunc(getFileUrl, h.getFile)
	mux.HandleFunc(getFileMetadataUrl, h.getFileMetadata)
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
func (h *Handler) uploadFile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.writeErrorResponse(writer, errInvalidMethod)
		return
	}

	file, header, err := h.getFormFile(writer, request, "file")

	if err != nil {
		return
	}

	defer file.Close()

	input := &service.CreateMediaInput{
		Name:        header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Size:        header.Size,
		File:        file,
	}

	mediaId, err := h.services.Media.Create(request.Context(), input)

	if err != nil {
		h.writeDomainErrorResponse(writer, err)
		return
	}

	h.writeSuccessResponse(writer, mediaId)
}

// GetMedia godoc
// @Summary      Получить медиафайл
// @Tags         media
// @Accept       */*
// @Content-Type */*
// @param        id   path      string  true  "media ID"
// @Success      200
// @Failure      404
// @Router       /v1/media/{id} [get]
func (h *Handler) getFile(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	mediaId := strings.TrimPrefix(request.RequestURI, getFileUrl)

	if mediaId == "" {
		h.logger.Println("empty string")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	mediaData, err := h.services.Media.GetFile(request.Context(), mediaId)

	if err != nil {
		if err == domain.ErrMediaNotFound {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		h.logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer mediaData.File.Close()
	bytes := make([]byte, mediaData.Size)
	_, err = mediaData.File.Read(bytes)

	if err != nil {
		h.logger.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", mediaData.ContentType)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", mediaData.Size))
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(bytes)

	if err != nil {
		h.logger.Println(err)
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

// ChangePassword godoc
// @Summary      Получить метаинформацию о файле
// @Tags         media
// @Accept       */*
// @Produce      json
// @Content-Type application/json
// @param        id   path      string  true  "media ID"
// @Success      200 {object} apiResponse{result=fileMetadataResponse}
// @Router       /v1/media/metadata/{id} [get]
func (h *Handler) getFileMetadata(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.writeErrorResponse(writer, errInvalidMethod)
		return
	}

	mediaId := strings.TrimPrefix(request.RequestURI, getFileMetadataUrl)

	if mediaId == "" {
		h.logger.Println("empty string")
		h.writeErrorResponse(writer, errEmptyId)
		return
	}

	fileMetadata, err := h.services.Media.GetMetadata(request.Context(), mediaId)

	if err != nil {
		h.writeDomainErrorResponse(writer, err)
		return
	}

	response := fileMetadataResponse{
		ID:          fileMetadata.ID,
		Name:        fileMetadata.Name,
		Size:        fileMetadata.Size,
		ContentType: fileMetadata.ContentType,
	}

	h.writeSuccessResponse(writer, response)
}
