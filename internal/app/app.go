package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/app/transfer"
	"github.com/hassanalgoz/swe/internal/ent"
	"github.com/hassanalgoz/swe/internal/outbound/metrics"
)

type App struct {
	ctx      context.Context
	transfer transfer.Subdomain
}

func New(
	ctx context.Context,
) App {
	transfer := transfer.New()

	return App{
		ctx:      ctx,
		transfer: transfer,
	}
}

// MoneyTransfer moves money from one account to another
// errors:
// - ErrNotFound: either from- or to-account not found
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (a *App) MoneyTransfer(from, to uuid.UUID, amount int64) error {
	fromAccount, err := a.transfer.GetAccount(a.ctx, from)
	if err != nil {
		return err
	}
	toAccount, err := a.transfer.GetAccount(a.ctx, to)
	if err != nil {
		return err
	}

	err = a.transfer.ExecuteTransfer(a.ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	metrics.MyCounter.Inc()
	return nil
}

// GetAccount retrieves an account by id
// errors:
// - ErrNotFound: account not found
func (a *App) GetAccount(id uuid.UUID) (*ent.Account, error) {
	acc, err := a.transfer.GetAccount(a.ctx, id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
