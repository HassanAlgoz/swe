package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
<<<<<<< HEAD:internal/contexts/user/methods.go
	"github.com/hassanalgoz/swe/internal/entities"
	"github.com/hassanalgoz/swe/internal/outbound/search"
=======
	"github.com/hassanalgoz/swe/internal/common"
	"github.com/hassanalgoz/swe/internal/outbound/database"
>>>>>>> 0b0bd21ff2d9b049103ec87d28cc967e581034a4:internal/app/user/methods.go
)

type NameChange struct {
	UserProfile common.UserProfile
	NewUsername string
}

type DomainContext struct {
<<<<<<< HEAD:internal/contexts/user/methods.go
	DB        *sql.DB
	svcSearch search.SearchServiceClient
}

func NewContext(db *sql.DB, svcSearch search.SearchServiceClient) DomainContext {
	return DomainContext{db, svcSearch}
}

func (dc *DomainContext) GetUserProfile(ctx context.Context, id uuid.UUID) (*entities.UserProfile, error) {
	var profile *entities.UserProfile
	row := dc.DB.QueryRowContext(ctx, "SELECT id, username FROM user_profile WHERE id = ?", id)
=======
	db *sql.DB
}

func NewContext() DomainContext {
	return DomainContext{
		db: database.Get(),
	}
}

func (dc *DomainContext) GetUserProfile(ctx context.Context, id uuid.UUID) (common.UserProfile, error) {
	var profile common.UserProfile
	row := dc.db.QueryRowContext(ctx, "SELECT id, username FROM user_profile WHERE id = ?", id)
>>>>>>> 0b0bd21ff2d9b049103ec87d28cc967e581034a4:internal/app/user/methods.go
	err := row.Scan(&profile.ID, &profile.Username)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// ChangeName
// errors cases:
// - len(newUsername) < 3
// - new name matches current name
<<<<<<< HEAD:internal/contexts/user/methods.go
func (dc *DomainContext) ChangeName(ctx context.Context, profile entities.UserProfile, newUsername string) error {
=======
func (dc *DomainContext) ChangeName(ctx context.Context, profile common.UserProfile, newUsername string) error {
>>>>>>> 0b0bd21ff2d9b049103ec87d28cc967e581034a4:internal/app/user/methods.go
	// Validate the field itself

	if yes, msg := isValidUsername(newUsername); yes {
		return fmt.Errorf(`invalid username: "%s"`, msg)
	}

	// Validate field against persisted data
	if newUsername == profile.Username {
		return fmt.Errorf("new name matches current name")
	}

	tx, err := dc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res1, err := tx.ExecContext(ctx, "UPDATE user_profile SET username = ? WHERE id = ?", newUsername, profile.ID)
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

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (dc *DomainContext) Search(ctx context.Context, query string) ([]uuid.UUID, error) {
	resp, err := dc.svcSearch.GetSearchResults(ctx, &search.SearchRequest{
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	var result []uuid.UUID
	for _, r := range resp.GetResults() {
		id, err := uuid.Parse(r.GetItemId())
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}
