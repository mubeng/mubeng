package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(address string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Address() string {
	return s.httpServer.Addr
}
