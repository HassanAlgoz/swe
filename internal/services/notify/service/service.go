package service

import (
	"context"

	"github.com/hassanalgoz/swe/micro/lms/service/store"
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
		store: store.New(database.Get()),
	}
}
