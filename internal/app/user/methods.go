package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hassanalgoz/swe/internal/common"
)

type NameChange struct {
	UserProfile common.UserProfile
	NewUsername string
}

type DomainContext struct {
	DB *sql.DB
}

func NewContext(db *sql.DB) DomainContext {
	return DomainContext{db}
}

func (dc *DomainContext) GetUserProfile(ctx context.Context, id uuid.UUID) (common.UserProfile, error) {
	var profile common.UserProfile
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
func (dc *DomainContext) ChangeName(ctx context.Context, profile common.UserProfile, newUsername string) error {
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
