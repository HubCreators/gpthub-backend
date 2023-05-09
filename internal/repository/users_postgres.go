package repository

import (
	"auth/internal/domain"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UsersPostgres struct {
	db *pgxpool.Pool
}

func NewUsersPostgres(db *pgxpool.Pool) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (r *UsersPostgres) Create(ctx context.Context, user domain.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, email, password) VALUES ($1, $2, $3) RETURNING id", usersCollections)
	row := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersPostgres) GetByCredentials(ctx context.Context, email string, password string) (domain.User, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password = $2", usersCollections)
	row := r.db.QueryRow(ctx, query, email, password)
	logrus.Info(query, email, password)
	var user domain.User
	err := row.Scan(&user.ID)
	return user, err
}

func (r *UsersPostgres) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email FROM %s WHERE (session).\"refreshToken\" = $1", usersCollections)
	row := r.db.QueryRow(ctx, query, refreshToken)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
	return user, err
}

func (r *UsersPostgres) SetSession(ctx context.Context, userID string, session domain.Session) error {
	query := fmt.Sprintf("UPDATE %s SET session = ROW($1, $2) WHERE id = $3", usersCollections)
	logrus.Info(query)
	_, err := r.db.Exec(ctx, query, session.RefreshToken, session.ExpiresAt, userID)
	return err
}
