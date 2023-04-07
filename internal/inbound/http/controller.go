package http

import (
	"context"
	"log"
	"net/http"

	"github.com/hassanalgoz/swe/internal/actions"
)

type Controller struct {
	ctx     context.Context
	mux     *http.ServeMux
	actions actions.Actions
}

func NewController(ctx context.Context, acts actions.Actions) *Controller {
	mux := http.NewServeMux()
	c := &Controller{
		mux:     mux,
		ctx:     ctx,
		actions: acts,
	}
	c.registerHandlers()
	return c
}

// Listen calls http.ListenAndServe
func (c *Controller) Listen(addr string) error {
	log.Println("Server listening on port 8080...")
	return http.ListenAndServe(addr, c.mux)
}
