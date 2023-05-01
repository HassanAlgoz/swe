package transfer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/common"
	"github.com/hassanalgoz/swe/internal/outbound/database"
	"github.com/hassanalgoz/swe/internal/outbound/grpc/search"
)

type Subdomain struct {
	db     *sql.DB
	search *search.Client
}

func New() Subdomain {
	return Subdomain{
		db:     database.Get(),
		search: search.Get(),
	}
}

// GetAccount
// errors:
// - ErrNotFound: account with given id is not found
func (s *Subdomain) GetAccount(ctx context.Context, id uuid.UUID) (*common.Account, error) {
	account := &common.Account{}
	row := s.db.QueryRowContext(ctx, "SELECT id, name, email, currency, freezed_since FROM accounts WHERE id = ?", id)
	err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Currency, &account.FreezedSince)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account id: %w", common.ErrNotFound)
		}
		return nil, err
	}
	return account, nil
}

// SaveTransfer
// errors:
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (s *Subdomain) SaveTransfer(ctx context.Context, from, to *common.Account, amount int64) error {
	// Validate the field itself
	if amount <= 0 {
		return &common.ErrInvalidArgument{
			Argument: "amount",
			Message:  fmt.Sprintf("expected: amount <= 0, got: %v", amount),
		}
	}

	// Validate the state of persisted data
	if yes, msg := isFreezed(from); yes {
		return &common.ErrInvalidState{
			RelatedArgument: "from",
			Message:         msg,
		}
	}

	if yes, msg := isFreezed(from); yes {
		return &common.ErrInvalidState{
			RelatedArgument: "to",
			Message:         msg,
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
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
