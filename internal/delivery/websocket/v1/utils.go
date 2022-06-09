package v1

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
)

const idRegexp = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

type userData struct {
	EventID     string
	UserID      string
	UserName    string
	UserImageID string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) upgradeConnection(writer http.ResponseWriter, request *http.Request) (userData, *websocket.Conn, error) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return userData{}, nil, errors.New("method not allowed")
	}

	conn, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return userData{}, nil, err
	}

	paramToken := request.URL.Query().Get("accessToken")

	if paramToken == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return userData{}, nil, errors.New("empty accessToken")
	}

	token, err := h.tokenManager.ParseToken(paramToken)

	if err != nil {
		if err == auth.ErrInvalidToken {
			writer.WriteHeader(http.StatusUnauthorized)
			return userData{}, nil, err
		}

		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return userData{}, nil, err
	}

	isValid, err := validateId(token.ID)

	if err != nil {
		h.logger.LogError(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return userData{}, nil, err
	}

	if !isValid {
		writer.WriteHeader(http.StatusBadRequest)
		return userData{}, nil, errors.New("invalid id")
	}

	return userData{
		UserID:      token.ID,
		UserName:    token.Name,
		UserImageID: token.ImageID,
	}, conn, nil
}

type readContext struct {
	context    context.Context
	conn       *websocket.Conn
	subscriber service.Subscriber
	userData   userData
}

func readMessages(context *readContext,
	messageHandler func(ctx context.Context, userData userData, body []byte),
	closeHandler func(ctx context.Context, userData userData, subscriber service.Subscriber)) {
	defer closeHandler(context.context, context.userData, context.subscriber)
	for {
		messageType, data, err := context.conn.ReadMessage()

		if err != nil {
			return
		}

		switch messageType {
		case websocket.TextMessage:
			messageHandler(context.context, context.userData, data)
		case websocket.CloseMessage:
			_ = context.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"))
			return
		case websocket.PingMessage:
			err = context.conn.WriteMessage(websocket.PongMessage, data)
			if err != nil {
				return
			}
		default:
			_ = context.conn.WriteMessage(websocket.TextMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "closing connection"))
			return
		}
	}
}

func writeMessages(pingPeriod time.Duration, conn *websocket.Conn, ch <-chan []byte) {
	ticker := time.NewTicker(pingPeriod * 9 / 10)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	err := conn.SetReadDeadline(time.Now().Add(pingPeriod))

	if err != nil {
		return
	}

	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pingPeriod))
	})

	for {
		select {
		case <-ticker.C:
			if err = conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case message, ok := <-ch:
			if !ok {
				ch = nil
				return
			}
			if err = conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
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
