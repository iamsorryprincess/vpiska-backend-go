package v1

import (
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

const idRegexp = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

type Handler struct {
	logger       logger.Logger
	tokenManager auth.TokenManager
	events       service.Events
	publisher    service.Publisher
}

func NewHandler(logger logger.Logger, tokenManager auth.TokenManager,
	events service.Events,
	publisher service.Publisher) *Handler {
	return &Handler{
		logger:       logger,
		tokenManager: tokenManager,
		events:       events,
		publisher:    publisher,
	}
}

func (h *Handler) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if strings.HasSuffix(request.URL.Path, "/event") {
			h.upgradeEventConnection(writer, request)
			return
		}

		writer.WriteHeader(http.StatusNotFound)
	})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) upgradeConnection(writer http.ResponseWriter, request *http.Request) (socketContext, *connection, error) {
	conn, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return socketContext{}, nil, err
	}

	params := strings.Split(request.URL.RawQuery, "&")
	paramToken := getQueryValue("accessToken", params)

	if paramToken == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return socketContext{}, nil, errors.New("empty accessToken")
	}

	token, err := h.tokenManager.ParseToken(paramToken)

	if err != nil {
		if err == auth.ErrInvalidToken {
			writer.WriteHeader(http.StatusUnauthorized)
			return socketContext{}, nil, err
		}

		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return socketContext{}, nil, err
	}

	isValid, err := validateId(token.ID)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return socketContext{}, nil, err
	}

	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		return socketContext{}, nil, errors.New("invalid id")
	}

	ctx := socketContext{
		Context:     request.Context(),
		UserID:      token.ID,
		UserName:    token.Name,
		UserImageID: token.ImageID,
	}

	socketConnection := &connection{
		mutex: sync.Mutex{},
		Conn:  conn,
	}

	return ctx, socketConnection, nil
}

func (h *Handler) readMessages(ctx socketContext, conn *connection,
	subscriber service.Subscriber,
	messageHandler func(ctx socketContext, body []byte),
	closeHandler func(ctx socketContext, subscriber service.Subscriber)) {
	defer conn.Close()
	defer closeHandler(ctx, subscriber)
	for {
		messageType, data, err := conn.ReadMessage()

		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				return
			}
			var netOpErr *net.OpError
			if errors.As(err, &netOpErr) && netOpErr.Op == "read" {
				return
			}
			h.logger.LogError(err)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			messageHandler(ctx, data)
		case websocket.CloseMessage:
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"))
			if err != nil {
				h.logger.LogError(err)
			}
			return
		case websocket.PingMessage:
			err = conn.WriteMessage(websocket.PongMessage, data)
			if err != nil {
				h.logger.LogError(err)
				return
			}
		default:
			err = conn.WriteMessage(websocket.TextMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "closing connection"))
			if err != nil {
				h.logger.LogError(err)
			}
			return
		}
	}
}

func validateId(id string) (bool, error) {
	if id == "" {
		return false, nil
	}

	if matched, err := regexp.MatchString(idRegexp, id); err != nil {
		return false, err
	} else if !matched {
		return false, nil
	}

	return true, nil
}

func getQueryValue(name string, splitQueryStr []string) string {
	for _, param := range splitQueryStr {
		if strings.HasPrefix(param, name) {
			return strings.TrimPrefix(param, name+"=")
		}
	}
	return ""
}
