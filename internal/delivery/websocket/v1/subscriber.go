package v1

import (
	"github.com/gorilla/websocket"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
	"github.com/iamsorryprincess/vpiska-backend-go/pkg/logger"
)

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

func (s *subscriber) OnClose() {
	err := s.conn.Close()

	if err != nil {
		s.logger.LogError(err)
	}
}
