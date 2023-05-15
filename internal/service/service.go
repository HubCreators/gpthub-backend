package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/auth"
	"auth/pkg/hash"
	"context"
	openai "github.com/sashabaranov/go-openai"
	"time"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	SignUp(ctx context.Context, username string, email string, password string) (int, error)
	SignIn(ctx context.Context, email string, password string) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type OpenAI interface {
	Communicate(content string) (string, error)
}

type Services struct {
	Users
	OpenAI
}

type Deps struct {
	Repos           *repository.Repository
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration

	OpenAIClient *openai.Client
	OpenAIModel  string
}

func NewService(deps Deps) *Services {
	userService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)
	openAIService := NewOpenAIService(deps.OpenAIClient, deps.OpenAIModel)

	return &Services{
		Users:  userService,
		OpenAI: openAIService,
	}
}
