package service

import (
	"fmt"
	"net/http"

	"github.com/RacoonMediaServer/rms-mirror/internal/config"
)

type Service struct {
	cfg config.Configuration
	s   http.Server
}

func New(cfg config.Configuration) *Service {
	return &Service{
		cfg: cfg,
		s: http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port),
		},
	}
}

func (s *Service) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.proxyFunc)

	s.s.Handler = mux
	return s.s.ListenAndServe()
}
