package v1

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/auth"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

func (h *Handler) initWebSockets(router *gin.RouterGroup) {
	websockets := router.Group("websockets")
	websockets.GET("/event", h.upgradeEventConnection)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) upgradeEventConnection(context *gin.Context) {
	eventId := context.Query("eventId")
	validationErrs, err := validateId(eventId)

	if err != nil {
		h.logger.LogError(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	if len(validationErrs) > 0 {
		context.Status(http.StatusBadRequest)
		return
	}

	ctx, conn, err := h.upgradeConnection(context)

	if err != nil {
		return
	}

	ctx.EventID = eventId
	sub := newSubscriber(h.logger, conn)
	h.services.Publisher.Subscribe(eventId, sub)
	err = h.services.Events.AddUserInfo(context.Request.Context(), service.AddUserInfoInput{
		EventID: eventId,
		UserID:  ctx.UserID,
	})

	if err != nil {
		h.services.Publisher.Unsubscribe(eventId, sub)
		closeErr := conn.Close()
		if closeErr != nil {
			h.logger.LogError(closeErr)
		}
		if domain.IsInternalError(err) {
			h.logger.LogError(err)
			context.Status(http.StatusInternalServerError)
			return
		}
		context.Status(http.StatusBadRequest)
		return
	}

	go h.readMessages(ctx, conn, sub, h.eventHandler, h.eventCloseHandler)
}

func (h *Handler) upgradeRangeConnection(context *gin.Context) {
	ctx, conn, err := h.upgradeConnection(context)

	if err != nil {
		return
	}

	sub := newSubscriber(h.logger, conn)
	go h.readMessages(ctx, conn, sub, h.rangeHandler, h.rangeCloseHandler)
}

func (h *Handler) upgradeConnection(context *gin.Context) (socketContext, *connection, error) {
	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		h.logger.LogError(err)
		context.Status(http.StatusInternalServerError)
		return socketContext{}, nil, err
	}

	token, err := h.tokenManager.ParseToken(context.Query("accessToken"))

	if err != nil {
		if err == auth.ErrInvalidToken {
			context.Status(http.StatusUnauthorized)
			return socketContext{}, nil, err
		}

		h.logger.LogError(err)
		context.Status(http.StatusInternalServerError)
		return socketContext{}, nil, err
	}

	validationErrs, err := validateId(token.ID)

	if err != nil {
		h.logger.LogError(err)
		context.Status(http.StatusInternalServerError)
		return socketContext{}, nil, err
	}

	if len(validationErrs) > 0 {
		context.Status(http.StatusBadRequest)
		return socketContext{}, nil, errInvalidId
	}

	ctx := socketContext{
		Context:     context.Request.Context(),
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
			h.logger.LogError(err)
			err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "closing connection"))
			if err != nil {
				h.logger.LogError(err)
			}
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

type socketContext struct {
	EventID     string
	UserID      string
	UserName    string
	UserImageID string
	context.Context
	ch <-chan struct{}
}

func (c socketContext) Done() <-chan struct{} {
	return c.ch
}

func (c socketContext) Err() error {
	select {
	case <-c.ch:
		return context.Canceled
	default:
		return nil
	}
}

type connection struct {
	mutex sync.Mutex
	*websocket.Conn
}

func (c *connection) Close() error {
	c.mutex.Lock()
	err := c.Conn.Close()
	c.mutex.Unlock()
	return err
}

func (c *connection) WriteMessage(messageType int, data []byte) error {
	c.mutex.Lock()
	err := c.Conn.WriteMessage(messageType, data)
	c.mutex.Unlock()
	return err
}

type subscriber struct {
	logger logger.Logger
	conn   *connection
}

func newSubscriber(logger logger.Logger, connection *connection) service.Subscriber {
	return &subscriber{
		logger: logger,
		conn:   connection,
	}
}

func (s *subscriber) OnReceive(message []byte) {
	err := s.conn.WriteMessage(websocket.TextMessage, message)

	if err != nil {
		s.logger.LogError(err)
	}
}
