package http

import (
	"context"
	"log"
	"net/http"

	"github.com/hassanalgoz/swe/internal/app"
	"github.com/hassanalgoz/swe/internal/outbound/logger"
	"github.com/rs/zerolog"
)

type Server struct {
	ctx    context.Context
	app    app.App
	mux    *http.ServeMux
	logger zerolog.Logger
}

func NewServer(ctx context.Context, a app.App) *Server {
	mux := http.NewServeMux()
	c := &Server{
		mux:    mux,
		ctx:    ctx,
		app:    a,
		logger: logger.Get(),
	}
	c.registerHandlers()
	return c
}

// Listen calls http.ListenAndServe
func (s *Server) Listen(addr string) error {
	log.Printf("Server listening on port %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
