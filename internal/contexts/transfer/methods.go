package transfer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/entities"
)

type DomainContext struct {
	db *sql.DB
}

func NewContext(db *sql.DB) DomainContext {
	return DomainContext{
		db: db,
	}
}

// GetAccount
// errors:
// - ErrNotFound: account with given id is not found
func (dc *DomainContext) GetAccount(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	account := &entities.Account{}
	row := dc.db.QueryRowContext(ctx, "SELECT id, name, email, currency, freezed_since FROM accounts WHERE id = ?", id)
	err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Currency, &account.FreezedSince)
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
func (dc *DomainContext) SaveTransfer(ctx context.Context, from, to *entities.Account, amount int64) error {
	// Validate the field itself
	if amount <= 0 {
		return &entities.ErrInvalidArgument{
			Argument: "amount",
			Message:  fmt.Sprintf("expected: amount <= 0, got: %v", amount),
		}
	}

	// Validate the state of persisted data
	yes, msg := isFreezed(from)
	if yes {
		return &entities.ErrInvalidState{
			RelatedArgument: "from",
			Message:         msg,
		}
	}

	yes, msg = isFreezed(from)
	if yes {
		return &entities.ErrInvalidState{
			RelatedArgument: "to",
			Message:         msg,
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
