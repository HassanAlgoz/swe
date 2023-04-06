package actions

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/contexts/transfer"
)

type Actions struct {
	ctx context.Context
	db  *sql.DB
}

func New(ctx context.Context, db *sql.DB) Actions {
	return Actions{
		ctx: ctx,
		db:  db,
	}
}

// MoneyTransfer moves money from one account to another
// errors:
// - ErrNotFound: either from- or to-account not found
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (a *Actions) MoneyTransfer(from, to uuid.UUID, amount int64) error {
	transferContext := transfer.NewContext(a.ctx, a.db)

	fromAccount, err := transferContext.GetAccount(a.ctx, from)
	if err != nil {
		return err
	}
	toAccount, err := transferContext.GetAccount(a.ctx, to)
	if err != nil {
		return err
	}

	err = transferContext.SaveTransfer(a.ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	return nil
}
