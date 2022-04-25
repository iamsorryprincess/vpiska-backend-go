package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/commands"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"
)

type UsersHandler struct {
	createCommandHandler *commands.CreateUserHandler
}

func InitUsersHandler(createUserHandler *commands.CreateUserHandler) *UsersHandler {
	return &UsersHandler{
		createCommandHandler: createUserHandler,
	}
}

func InitV1UsersHandler(group *gin.RouterGroup, usersHandler *UsersHandler) {
	v1 := group.Group("/users")
	v1.POST("/create", createUserHandler(usersHandler.createCommandHandler))
}

func createUserHandler(commandHandler *commands.CreateUserHandler) gin.HandlerFunc {
	return func(context *gin.Context) {
		command := &commands.CreateUserCommand{}
		parseError := context.BindJSON(command)

		if parseError != nil {
			fmt.Println(parseError)
			context.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := commandHandler.Handle(command)

		if err != nil {
			response := createErrorResponse[models.UserResponse](err)
			context.JSON(http.StatusOK, response)
			return
		}

		response := createSuccessResponse(result)
		context.JSON(http.StatusOK, response)
	}
}
