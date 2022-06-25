package server

import (
	"context"
	"fmt"
	"net/http"
)

type HttpServer struct {
	httpServer *http.Server
}

func NewHttpServer(port int, handler http.Handler) *HttpServer {
	return &HttpServer{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: handler,
		},
	}
}

func (s *HttpServer) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
