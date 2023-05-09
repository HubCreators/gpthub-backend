package service

import (
	"auth/internal/repository"
	"auth/pkg/auth"
	"auth/pkg/hash"
	"context"
	"time"
)

type UserSignUpInput struct {
	Username string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, inputUser UserSignUpInput) (int, error)
	SignIn(ctx context.Context, inputUser UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Services struct {
	Users
}

type Deps struct {
	Repos           *repository.Repository
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewService(deps Deps) *Services {
	userService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users: userService,
	}
}
