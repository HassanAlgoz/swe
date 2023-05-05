package service

import (
	"context"
	"net/http"

	"github.com/hassanalgoz/swe/pkg/adapters/services/lms"
	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/rs/zerolog"
)

type service struct {
	ctx    context.Context
	mux    *http.ServeMux
	logger zerolog.Logger
}

var lmsClient = lms.Singleton()
var log = logger.Singleton()

func NewServer(ctx context.Context) *service {
	c := &service{
		ctx:    ctx,
		mux:    http.NewServeMux(),
		logger: log,
	}
	c.registerHandlers()
	return c
}

func (s *service) Listen(addr string) error {
	log.Info().Msgf("Server listening on port %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
