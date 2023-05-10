package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/auth"
	"auth/pkg/hash"
	"context"
	"time"
)

type UsersService struct {
	repo         repository.Users
	hasher       hash.PasswordHasher
	tokenManager auth.TokenManager

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUsersService(repo repository.Users, hasher hash.PasswordHasher, tokenManager auth.TokenManager,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UsersService) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *UsersService) SignUp(ctx context.Context, username string, email string, password string) (int, error) {
	if s.isUserExists(ctx, email) {
		return 0, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return 0, err
	}

	user := domain.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIn(ctx context.Context, email string, password string) (Tokens, error) {
	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, email, hashedPassword)

	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, user.ID)
}

func (s *UsersService) RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return Tokens{}, err
	}
	return s.createSession(ctx, user.ID)
}

func (s *UsersService) createSession(ctx context.Context, userId string) (Tokens, error) {
	var (
		tokens Tokens
		err    error
	)

	tokens.AccessToken, err = s.tokenManager.NewJWT(userId, s.accessTokenTTL)
	if err != nil {
		return tokens, err
	}

	tokens.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return tokens, err
	}

	session := domain.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, userId, session)
	return tokens, err
}

func (s *UsersService) isUserExists(ctx context.Context, email string) bool {
	user, _ := s.GetByEmail(ctx, email)
	if user.Username != "" {
		return true
	}
	return false
}
