package v1

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
)

const idRegexp = `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) upgradeConnection(writer http.ResponseWriter, request *http.Request) (socketContext, *websocket.Conn, error) {
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

	return ctx, conn, nil
}

func readMessages(ctx socketContext, conn *websocket.Conn, subscriber service.Subscriber,
	messageHandler func(ctx socketContext, body []byte),
	closeHandler func(ctx socketContext, subscriber service.Subscriber)) {
	defer closeHandler(ctx, subscriber)
	for {
		messageType, data, err := conn.ReadMessage()

		if err != nil {
			return
		}

		switch messageType {
		case websocket.TextMessage:
			messageHandler(ctx, data)
		case websocket.CloseMessage:
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing connection"))
			return
		case websocket.PingMessage:
			err = conn.WriteMessage(websocket.PongMessage, data)
			if err != nil {
				return
			}
		default:
			_ = conn.WriteMessage(websocket.TextMessage, websocket.FormatCloseMessage(websocket.CloseUnsupportedData, "closing connection"))
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

func getQueryValue(name string, splitQueryStr []string) string {
	for _, param := range splitQueryStr {
		if strings.HasPrefix(param, name) {
			return strings.TrimPrefix(param, name+"=")
		}
	}
	return ""
}
