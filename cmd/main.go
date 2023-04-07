package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hassanalgoz/swe/internal/actions"
	"github.com/hassanalgoz/swe/internal/inbound/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, os.Interrupt)
	go func() {
		<-cancelChan
		cancel()
	}()

	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/bank")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	act := actions.New(ctx, db)
	httpController := http.NewController(ctx, act)
	if err = httpController.Listen(":8080"); err != nil {
		panic(err)
	}
}
