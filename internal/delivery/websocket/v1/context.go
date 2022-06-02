package v1

import "context"

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
