package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/actions"
)

// Response conforms to: https://google.github.io/styleguide/jsoncstyleguide.xml
type Response struct {
	Data  any           `json:"data,omitempty"`
	Error ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Errors  []Error `json:"errors"`
}

type Error struct {
	Location     string `json:"location"`      // "Authorization" | "userId" | ...etc
	LocationType string `json:"location_type"` // "header" | "parameter"
	Message      string `json:"message"`
	Reason       string `json:"reason"`
}

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

	mux.HandleFunc("/actions:transfer-money", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(Response{
				Error: ResponseError{
					Code:    http.StatusMethodNotAllowed,
					Message: "Method Not Allowed",
				},
			})
			return
		}

		// Read request headers, parse, validate, ...etc.
		userID := r.Header.Get("x-user-id")
		// Parse, validate.....
		if userID == "" {
			userID = "default-user-id"
		}

		// Apply rules based on headers...
		// ...
		// ...

		// Parse Body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: ResponseError{
					Code:    http.StatusBadRequest,
					Message: "Invalid request body",
				},
			})
			return
		}

		// Parse JSON
		var fields struct {
			From   string `json:"from"`
			To     string `json:"to"`
			Amount int64  `json:"amount"`
		}
		err = json.Unmarshal(body, &fields)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: ResponseError{
					Code:    http.StatusBadRequest,
					Message: "Invalid json",
				},
			})
			return
		}

		// Parse fields
		var errors []Error
		from, err := uuid.Parse(fields.From)
		if err != nil {
			errors = append(errors, Error{
				LocationType: "parameter",
				Location:     "from",
				Message:      "invalid uuid",
				Reason:       err.Error(),
			})
		}
		to, err := uuid.Parse(fields.To)
		if err != nil {
			errors = append(errors, Error{
				LocationType: "parameter",
				Location:     "to",
				Message:      "invalid uuid",
				Reason:       err.Error(),
			})
		}

		amount := fields.Amount

		// Check for bad parameters
		if len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: ResponseError{
					Code:    http.StatusBadRequest,
					Message: "Invalid parameters",
					Errors:  errors,
				},
			})
			return
		}

		// Log before dispatch
		log.Printf("Transfer request from %s to %s for %d by user %s", from, to, amount, userID)

		// Dispatch
		err = actions.MoneyTransfer(ctx, from, to, amount)
		// switch err {
		// 	// TODO
		// }
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Error: ResponseError{
					Code:    http.StatusBadRequest,
					Message: "Invalid json",
				},
			})
			return
		}

		// Done
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(Response{
			Data: nil,
		})
	})

	fmt.Println("Server listening on port 8080...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
