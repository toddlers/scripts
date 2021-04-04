package account

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
)

var RepoErr = errors.New("Unable to handle Repo Request")

type repo struct {
	db     *sql.DB
	logger log.Logger
}

func NewRepo(db *sql.DB, logger log.Logger) Repository {
	return &repo{
		db:     db,
		logger: log.With(logger, "repo", "sql"),
	}
}

func (r repo) CreateUser(ctx context.Context, user User) error {
	sql := `
	INSERT INTO users (id, email, password)
	VALUES ($1, $2, $3)
	`
	if user.Email == "" || user.Password == "" {
		return RepoErr
	}
	_, err := r.db.ExecContext(ctx, sql,
		user.ID,
		user.Email,
		user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r repo) GetUser(ctx context.Context, id string) (string, error) {
	var email string
	err := r.db.QueryRow("SELECT email FROM users WHERE id=$1", id).Scan(&email)
	if err != nil {
		return "", RepoErr
	}
	return email, nil
}
