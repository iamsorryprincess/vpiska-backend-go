package v1

import (
	"github.com/iamsorryprincess/vpiska-backend-go/internal/service"
)

type subscriber struct {
	ch chan<- []byte
}

func newSubscriber(ch chan<- []byte) service.Subscriber {
	return &subscriber{
		ch: ch,
	}
}

func (s *subscriber) OnReceive(message []byte) {
	s.ch <- message
}

func (s *subscriber) OnClose() {
	close(s.ch)
}
