package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/entities"
)

type NameChange struct {
	UserProfile entities.UserProfile
	NewUsername string
}

type domainContext struct {
	DB *sql.DB
}

func NewContext(db *sql.DB) domainContext {
	return domainContext{db}
}

func (dc *domainContext) GetUserProfile(ctx context.Context, id uuid.UUID) (entities.UserProfile, error) {
	var profile entities.UserProfile
	row := dc.DB.QueryRowContext(ctx, "SELECT id, username FROM user_profile WHERE id = ?", id)
	err := row.Scan(&profile.ID, &profile.Username)
	if err != nil {
		return profile, err
	}
	return profile, nil
}

// ChangeName
// errors cases:
// - len(newUsername) < 3
// - new name matches current name
func (dc *domainContext) ChangeName(ctx context.Context, profile entities.UserProfile, newUsername string) error {
	// Validate the field itself

	if yes, msg := isValidUsername(newUsername); yes {
		return fmt.Errorf(`invalid username: "%s"`, msg)
	}

	// Validate field against persisted data
	if newUsername == profile.Username {
		return fmt.Errorf("new name matches current name")
	}

	tx, err := dc.DB.BeginTx(ctx, nil)
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
