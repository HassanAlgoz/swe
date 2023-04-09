package actions

import (
	"context"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/contexts/transfer"
	"github.com/hassanalgoz/swe/internal/contexts/user"
	"github.com/hassanalgoz/swe/internal/entities"
)

type Actions struct {
	ctx             context.Context
	transferContext transfer.DomainContext
	userContext     user.DomainContext
}

func New(
	ctx context.Context,
	transferContext transfer.DomainContext,
	userContext user.DomainContext,
) Actions {
	return Actions{
		ctx:             ctx,
		transferContext: transferContext,
		userContext:     userContext,
	}
}

// MoneyTransfer moves money from one account to another
// errors:
// - ErrNotFound: either from- or to-account not found
// - ErrInvalidArgument: amount <= 0
// - ErrInvalidState: either from- or to-account is freezed
func (a *Actions) MoneyTransfer(from, to uuid.UUID, amount int64) error {
	fromAccount, err := a.transferContext.GetAccount(a.ctx, from)
	if err != nil {
		return err
	}
	toAccount, err := a.transferContext.GetAccount(a.ctx, to)
	if err != nil {
		return err
	}

	err = a.transferContext.SaveTransfer(a.ctx, fromAccount, toAccount, amount)
	if err != nil {
		return err
	}
	return nil
}

// GetAccount retrieves an account by id
// errors:
// - ErrNotFound: account not found
func (a *Actions) GetAccount(id uuid.UUID) (*entities.Account, error) {
	acc, err := a.transferContext.GetAccount(a.ctx, id)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (a *Actions) GetUsersProfilesByQuery(query string) ([]*entities.UserProfile, error) {
	ids, err := a.userContext.Search(a.ctx, query)
	if err != nil {
		return nil, err
	}
	var result []*entities.UserProfile
	for i := range ids {
		p, err := a.userContext.GetUserProfile(a.ctx, ids[i])
		if err != nil {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}
