package entities

import "github.com/google/uuid"

type CurrencyCode string

const (
	CurrencySaudiRiyal CurrencyCode = "682"
)

type Account struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Currency CurrencyCode
	Freezed  bool
}
