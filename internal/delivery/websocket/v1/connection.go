package v1

import (
	"sync"

	"github.com/gorilla/websocket"
)

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
