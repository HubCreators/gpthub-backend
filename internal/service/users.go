package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/auth"
	"auth/pkg/hash"
	"context"
	"github.com/sirupsen/logrus"
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

func (s *UsersService) SignUp(ctx context.Context, inputUser UserSignUpInput) (int, error) {
	hashedPassword, err := s.hasher.Hash(inputUser.Password)
	if err != nil {
		logrus.Fatal("Unable to hash password")
		return 0, err
	}

	user := domain.User{
		Username: inputUser.Username,
		Email:    inputUser.Email,
		Password: hashedPassword,
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIn(ctx context.Context, inputUser UserSignInInput) (Tokens, error) {
	hashedPassword, err := s.hasher.Hash(inputUser.Password)
	if err != nil {
		logrus.Fatal("Unable to hash password")
		return Tokens{}, err
	}

	user, err := s.repo.GetByCredentials(ctx, inputUser.Email, hashedPassword)

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

	logrus.Info(s)
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
	if err != nil {
		logrus.Fatal("RED FLAG!!!")
	}
	return tokens, err
}
