package whoop

import (
	"database/sql"
	"errors"
	"time"
)

var ErrStateNotFound = errors.New("state not found or expired")

func StoreState(db *sql.DB, state string, userID *int, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)

	// Convert optional userID to a value that can be inserted (nil or int)
	var uid interface{}
	if userID == nil {
		uid = nil
	} else {
		uid = *userID
	}

	_, err := db.Exec(`
		INSERT INTO oauth_states (state, user_id, expires_at)
		VALUES (?, ?, ?)
		`, state, uid, expiresAt)
	return err
}

func ValidateAndDeleteState(db *sql.DB, state string) (int, error) {
	// Start a transaction so we can select FOR UPDATE and delete safely
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	var rowID int
	var expiresAt time.Time
	var userID sql.NullInt64

	// select the row for update (prevents race conditions where two request try to validate at the same time)
	row := tx.QueryRow(`
		SELECT id, user_id, expires_at
		FROM oauth_states
		WHERE state = ?
		FOR UPDATE
		`, state)

	if err := row.Scan(&rowID, &userID, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrStateNotFound
		}
		return 0, err
	}

	// If expired, delete and return not found
	if time.Now().After(expiresAt) {
		if _, err := tx.Exec(`DELETE FROM oauth_states WHERE id = ?`, rowID); err != nil {
			return 0, err
		}
		if err := tx.Commit(); err != nil {
			return 0, err
		}
		return 0, ErrStateNotFound
	}

	// Otherwise, delete state
	if _, err := tx.Exec(`DELETE FROM oauth_states WHERE id = ?`, rowID); err != nil {
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Return the associated user id (0 if not set)
	if userID.Valid {
		return int(userID.Int64), nil
	}
	return 0, nil
}

const defaultStateTTL = 5 * time.Minute
