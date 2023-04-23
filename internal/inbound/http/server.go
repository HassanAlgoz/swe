package http

import (
	"context"
	"log"
	"net/http"

	"github.com/hassanalgoz/swe/internal/app"
)

type Server struct {
	ctx context.Context
	app app.App
	mux *http.ServeMux
}

func NewServer(ctx context.Context, a app.App) *Server {
	mux := http.NewServeMux()
	c := &Server{
		mux: mux,
		ctx: ctx,
		app: a,
	}
	c.registerHandlers()
	return c
}

// Listen calls http.ListenAndServe
func (c *Server) Listen(addr string) error {
	log.Printf("Server listening on port %s", addr)
	return http.ListenAndServe(addr, c.mux)
}
