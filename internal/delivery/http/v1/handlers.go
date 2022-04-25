package v1

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	usersHandler *UsersHandler
}

func InitHandler(usersHandler *UsersHandler) *Handler {
	return &Handler{
		usersHandler: usersHandler,
	}
}

func InitV1ApiHandlers(r *gin.Engine, handler *Handler) {
	v1 := r.Group("/api/v1")
	InitV1UsersHandler(v1, handler.usersHandler)
}
