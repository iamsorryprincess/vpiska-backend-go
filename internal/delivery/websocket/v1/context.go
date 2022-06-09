package v1

import "context"

type socketContext struct {
	context.Context
	ch <-chan struct{}
}

func newSocketContext(ctx context.Context) context.Context {
	return socketContext{
		Context: ctx,
	}
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
