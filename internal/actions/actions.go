package actions

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/contexts/transfer"
)

// MoneyTransfer goes from one account to another
func MoneyTransfer(ctx context.Context, from, to uuid.UUID, amount int64) error {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/bank")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	transferContext := transfer.NewContext(db)

	fromAccount, err := transferContext.GetAccount(ctx, from)
	if err != nil {
		log.Fatal(err)
	}
	toAccount, err := transferContext.GetAccount(ctx, to)
	if err != nil {
		log.Fatal(err)
	}

	err = transferContext.SaveTransfer(ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	return nil
}
