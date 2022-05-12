package v1

import (
	"net/http"
	"strings"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

const multipart = "multipart/form-data"

func (h *Handler) initMediaAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/media", h.routeMediaHandle)
}

func (h *Handler) routeMediaHandle(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		h.uploadFile(writer, request)
		return
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
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
	if !strings.Contains(request.Header.Get("Content-Type"), multipart) {
		writer.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	err := request.ParseMultipartForm(10 << 20)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
		return
	}

	if request.MultipartForm.File["file"] == nil {
		h.writeDomainErrorResponse(writer, domain.ErrEmptyMedia)
		return
	}

	file, header, err := request.FormFile("file")

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		h.logger.Println(err)
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
