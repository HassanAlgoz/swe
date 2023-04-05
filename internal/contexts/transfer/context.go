package transfer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/entities"
)

type domainContext struct {
	DB *sql.DB
}

func NewContext(db *sql.DB) domainContext {
	return domainContext{db}
}

func (dc *domainContext) GetAccount(ctx context.Context, id uuid.UUID) (entities.Account, error) {
	var account entities.Account
	row := dc.DB.QueryRowContext(ctx, "SELECT id, name, email, currency, freezed FROM accounts WHERE id = ?", id)
	err := row.Scan(&account.ID, &account.Name, &account.Email, &account.Currency, &account.Freezed)
	if err != nil {
		return account, err
	}
	return account, nil
}

// SaveTransfer
// errors cases:
// - amount <= 0
// - the from account is freezed
// - the to account is freezed
func (dc *domainContext) SaveTransfer(ctx context.Context, from, to entities.Account, amount int64) error {
	// Validate the field itself
	if amount <= 0 {
		return fmt.Errorf(`invalid amount: "%d"`, amount)
	}

	// Validate the state of persisted data
	if from.Freezed {
		return fmt.Errorf(`the from_account is freezed: "%s"`, from.ID)
	}
	if to.Freezed {
		return fmt.Errorf(`the to_account is freezed: "%s"`, from.ID)
	}

	tx, err := dc.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res1, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, from.ID)
	if err != nil {
		return err
	}
	rowsAffected1, err := res1.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected1 == 0 {
		return fmt.Errorf("no rows affected")
	}

	res2, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, to.ID)
	if err != nil {
		return err
	}
	rowsAffected2, err := res2.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected2 == 0 {
		return fmt.Errorf("no rows affected")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
