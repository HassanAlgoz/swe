package port

import (
	"context"
	"net/http"

	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/hassanalgoz/swe/pkg/services/adapters/lms"
)

type service struct {
	ctx context.Context
	mux *http.ServeMux
}

var lmsClient = lms.Singleton()
var log = logger.Get()

func NewServer(ctx context.Context) *service {
	c := &service{
		ctx: ctx,
		mux: http.NewServeMux(),
	}
	c.registerHandlers()
	return c
}

func (s *service) Listen(addr string) error {
	log.Info().Msgf("Server listening on port %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
