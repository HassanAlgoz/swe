package actions

import (
	"context"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/contexts/transfer"
)

type Actions struct {
	ctx             context.Context
	transferContext transfer.DomainContext
}

func New(
	ctx context.Context,
	transferContext transfer.DomainContext,
) Actions {
	return Actions{
		ctx:             ctx,
		transferContext: transferContext,
	}
}

// MoneyTransfer moves money from one account to another
// errors:
// - ErrNotFound: either from- or to-account not found
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (a *Actions) MoneyTransfer(from, to uuid.UUID, amount int64) error {
	fromAccount, err := a.transferContext.GetAccount(a.ctx, from)
	if err != nil {
		return err
	}
	toAccount, err := a.transferContext.GetAccount(a.ctx, to)
	if err != nil {
		return err
	}

	err = a.transferContext.SaveTransfer(a.ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	return nil
}
