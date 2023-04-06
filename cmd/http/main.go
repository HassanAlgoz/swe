package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hassanalgoz/swe/internal/actions"
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

	mux := http.NewServeMux()

	act := actions.New(ctx, db)

	registerHandlers(ctx, mux, act)

	fmt.Println("Server listening on port 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
