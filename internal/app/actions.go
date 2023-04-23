package app

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/app/transfer"
	"github.com/hassanalgoz/swe/internal/common"
)

type Actions struct {
	ctx      context.Context
	transfer transfer.DomainContext
}

func New(
	ctx context.Context,
	db *sql.DB,
) Actions {
	transfer := transfer.NewContext(db)
	return Actions{
		ctx:      ctx,
		transfer: transfer,
	}
}

// MoneyTransfer moves money from one account to another
// errors:
// - ErrNotFound: either from- or to-account not found
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (a *Actions) MoneyTransfer(from, to uuid.UUID, amount int64) error {
	fromAccount, err := a.transfer.GetAccount(a.ctx, from)
	if err != nil {
		return err
	}
	toAccount, err := a.transfer.GetAccount(a.ctx, to)
	if err != nil {
		return err
	}

	err = a.transfer.SaveTransfer(a.ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	return nil
}

// GetAccount retrieves an account by id
// errors:
// - ErrNotFound: account not found
func (a *Actions) GetAccount(id uuid.UUID) (*common.Account, error) {
	acc, err := a.transfer.GetAccount(a.ctx, id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
