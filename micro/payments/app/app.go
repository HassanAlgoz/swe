package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/micro/payments/app/payments"
	"github.com/hassanalgoz/swe/pkg/entities"
	"github.com/hassanalgoz/swe/pkg/infra/metrics"
)

type App struct {
	ctx      context.Context
	payments payments.Subdomain
}

func New(
	ctx context.Context,
) App {
	return App{
		ctx:      ctx,
		payments: payments.New(),
	}
}

// MoneyTransfer moves money from one account to another
// errors:
// - ErrNotFound: either from- or to-account not found
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (a *App) MoneyTransfer(from, to uuid.UUID, amount int64) error {
	fromAccount, err := a.payments.GetAccount(a.ctx, from)
	if err != nil {
		return err
	}
	toAccount, err := a.payments.GetAccount(a.ctx, to)
	if err != nil {
		return err
	}

	err = a.payments.ExecuteTransfer(a.ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	metrics.MyCounter.Inc()
	return nil
}

// GetAccount retrieves an account by id
// errors:
// - ErrNotFound: account not found
func (a *App) GetAccount(id uuid.UUID) (*entities.Account, error) {
	acc, err := a.payments.GetAccount(a.ctx, id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
