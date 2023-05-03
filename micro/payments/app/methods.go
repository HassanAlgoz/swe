package transfer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/ent"
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
func (s *Subdomain) GetAccount(ctx context.Context, id uuid.UUID) (*ent.Account, error) {
	account := &ent.Account{}
	row := s.db.QueryRowContext(ctx, "SELECT id, name, email, currency, freezed_since FROM accounts WHERE id = ?", id)
	err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Currency, &account.FreezedSince)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account id: %w", ent.ErrNotFound)
		}
		return nil, err
	}
	return account, nil
}

// ExecuteTransfer
// errors:
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (s *Subdomain) ExecuteTransfer(ctx context.Context, from, to *ent.Account, amount int64) error {
	// Validate the field itself
	if amount <= 0 {
		return &ent.ErrInvalidArgument{
			Argument: "amount",
			Message:  fmt.Sprintf("expected: amount <= 0, got: %v", amount),
		}
	}

	// Validate the state of persisted data
	if yes, msg := isFreezed(from); yes {
		return &ent.ErrInvalidState{
			RelatedArgument: "from",
			Message:         msg,
		}
	}

	if yes, msg := isFreezed(from); yes {
		return &ent.ErrInvalidState{
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

func (s *Subdomain) CreateRTP(ctx context.Context, rtp ent.RTP) (*ent.RTP, error) {
	rtp.ID = uuid.New()
	rtp.State = ent.RTPStateInitiated
	query := "INSERT INTO rtp (id, payer, payee, amount, created_at, due, state) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, rtp.ID, rtp.Payer, rtp.Payee, rtp.Amount, rtp.CreatedAt, rtp.Due, rtp.State)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected != 1 {
		return nil, fmt.Errorf("unexpected number of rows affected: %d", rowsAffected)
	}
	return &rtp, nil
}
