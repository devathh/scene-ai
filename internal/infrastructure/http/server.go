package httpserver

import (
	"context"
	"net/http"

	"github.com/devathh/scene-ai/internal/common/config"
)

type Server struct {
	cfg *config.Config
	srv *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		cfg: cfg,
		srv: &http.Server{
			Addr:         cfg.Server.Addr,
			Handler:      handler,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

func (s *Server) Start() error {
	if s.cfg.Server.TLS.Enable {
		return s.srv.ListenAndServeTLS(
			s.cfg.Server.TLS.ServerCertPath,
			s.cfg.Server.TLS.ServerKeyPath,
		)
	}

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
