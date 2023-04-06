package transfer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/entities"
)

type domainContext struct {
	db  *sql.DB
	ctx context.Context
}

func NewContext(ctx context.Context, db *sql.DB) domainContext {
	return domainContext{
		ctx: ctx,
		db:  db,
	}
}

// GetAccount
// errors:
// - ErrNotFound: account with given id is not found
func (dc *domainContext) GetAccount(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	account := &entities.Account{}
	row := dc.db.QueryRowContext(ctx, "SELECT id, name, email, currency, freezed FROM accounts WHERE id = ?", id)
	err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Currency, &account.Freezed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account id: %w", entities.ErrNotFound)
		}
		return nil, err
	}
	return account, nil
}

// SaveTransfer
// errors:
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (dc *domainContext) SaveTransfer(ctx context.Context, from, to *entities.Account, amount int64) error {
	// Validate the field itself
	if amount <= 0 {
		return &entities.ErrInvalidArgument{
			Argument: "amount",
			Message:  fmt.Sprintf("expected: amount <= 0, got: %v", amount),
		}
	}

	// Validate the state of persisted data
	if from.Freezed {
		return &entities.ErrInvalidState{
			RelatedArgument: "from",
			Message:         "account is freezed",
		}
	}
	if to.Freezed {
		return &entities.ErrInvalidState{
			RelatedArgument: "to",
			Message:         "account is freezed",
		}
	}

	tx, err := dc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, from.ID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, to.ID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
