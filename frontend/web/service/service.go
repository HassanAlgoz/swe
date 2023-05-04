package http

import (
	"context"
	"net/http"

	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/rs/zerolog"
)

type service struct {
	ctx    context.Context
	mux    *http.ServeMux
	logger zerolog.Logger
}

func NewServer(ctx context.Context) *service {
	c := &service{
		ctx:    ctx,
		mux:    http.NewServeMux(),
		logger: logger.Get(),
	}
	c.registerHandlers()
	return c
}

func (s *service) Listen(addr string) error {
	s.logger.Info().Msgf("Server listening on port %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
