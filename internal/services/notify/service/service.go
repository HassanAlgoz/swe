package service

import (
	"context"

	store "github.com/hassanalgoz/swe/internal/services/notify/store/port"
	"github.com/hassanalgoz/swe/pkg/infra/database"
)

type service struct {
	ctx   context.Context
	store store.Querier
}

func New(
	ctx context.Context,
) service {
	return service{
		ctx:   ctx,
		store: store.New(database.Get("notify")),
	}
}
