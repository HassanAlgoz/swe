package port

import (
	"context"
	"net/http"

	"github.com/hassanalgoz/swe/pkg/infra/logger"
	"github.com/hassanalgoz/swe/pkg/services/adapters/lms"
	lmsPort "github.com/hassanalgoz/swe/pkg/services/ports/lms"
)

var log = logger.Get()

type service struct {
	ctx context.Context
	mux *http.ServeMux
	lms lmsPort.LMSClient
}

func NewServer(ctx context.Context) *service {
	c := &service{
		ctx: ctx,
		mux: http.NewServeMux(),
		lms: lms.New(),
	}
	c.registerHandlers()
	return c
}

func (s *service) Listen(addr string) error {
	log.Info().Msgf("Server listening on port %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
