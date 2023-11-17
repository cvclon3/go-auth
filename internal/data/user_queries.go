// internal/data/user_queries.go
package data

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (um UserModel) Insert(user *User) (*UserID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := um.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var userID uuid.UUID

	query_user := `
    INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	args_user := []interface{}{user.Email, user.Password.hash, user.FirstName, user.LastName}

	if err := tx.QueryRowContext(ctx, query_user, args_user...).Scan(&userID); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}

	}

	query_user_profile := `
    INSERT INTO user_profile (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING RETURNING user_id`

	_, err = tx.ExecContext(ctx, query_user_profile, userID)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	id := UserID{
		Id: userID,
	}

	return &id, nil
}
