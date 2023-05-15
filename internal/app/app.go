package app

import (
	"auth/internal/config"
	"auth/internal/handler"
	"auth/internal/repository"
	"auth/internal/server"
	"auth/internal/service"
	"auth/pkg/auth"
	"auth/pkg/database/postgres"
	"auth/pkg/hash"
	"context"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

// @title           GPT API
// @version         1.0
// @description     rvinnie's ChatGPT server.

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey UsersAuth
// @in header
// @name Authorization

func Run(configPath string) {
	// Adding logger
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Initializing env variables
	if err := godotenv.Load(); err != nil {
		logrus.Error("Error loading .env file")
		return
	}

	// Initializing config
	cfg, err := config.Init(configPath)
	if err != nil {
		logrus.Error("Unable to parse config", err)
		return
	}

	// Initializing postgres
	db, err := postgres.NewConnPool(postgres.DBConfig{
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DBName:   cfg.Postgres.DBName,
	})
	if err != nil {
		logrus.Errorf("Unable to connect db: %v", err)
		return
	}
	defer db.Close()

	// Creating chat GPT client
	openAIClient := openai.NewClient(cfg.OpenAI.Token)
	openAIModel := openai.GPT3Dot5Turbo

	// Creating JWT token manager
	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		logrus.Errorf("Unable to create token manager: %v", err)
		return
	}

	// Creating hasher
	hasher := hash.NewSHA1Hasher(cfg.Auth.Salt)

	repository := repository.NewRepository(db)
	services := service.NewService(service.Deps{
		Repos:           repository,
		Hasher:          hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,

		OpenAIClient: openAIClient,
		OpenAIModel:  openAIModel,
	})
	handlers := handler.NewMainHandler(services, tokenManager)

	server := server.NewServer(cfg, handlers.InitRoutes(*cfg))
	go func() {
		server.Run()
	}()
	logrus.Info("AuthServer is running")

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // waiting SIGINT or SIGTERM

	logrus.Info("AuthServer shutting down")
	if err := server.Stop(context.Background()); err != nil {
		logrus.Errorf("Error on server shutting down: %s", err.Error())
	}
}
