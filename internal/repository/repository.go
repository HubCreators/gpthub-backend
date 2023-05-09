package repository

import (
	"auth/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users interface {
	Create(ctx context.Context, user domain.User) (int, error)
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	SetSession(ctx context.Context, userID string, session domain.Session) error
}

type Repository struct {
	Users
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Users: NewUsersPostgres(db),
	}
}
