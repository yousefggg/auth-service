package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/yousefggg/auth-service/internal/delivery"
	"github.com/yousefggg/auth-service/internal/repository/postgres"
	"github.com/yousefggg/auth-service/internal/usecase"
	"github.com/yousefggg/common-lib/pkg/jwt"
	"github.com/yousefggg/common-lib/pkg/logger"

	_ "github.com/yousefggg/auth-service/docs"
    "github.com/swaggo/http-swagger"
)
// @title           Auth Service API
// @version         1.0
// @description     Микросервис для аутентификации пользователей (регистрация и логин).
// @termsOfService  http://swagger.io/terms/

// @contact.name   Yousef Support
// @contact.url    http://github.com/yousefggg

// @host      localhost:8081
// @BasePath  /
// @accept    json
// @produce   json
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	logLevel := os.Getenv("ENV")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Init(logLevel)
	logger.Info("Starting auth-service", "env", logLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dsn := os.Getenv("DATABASE_URL")
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Failed to parse database config", "error", err)
		os.Exit(1)
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		logger.Error("Failed to create database pool", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	logger.Info("Successfully connected to database")

	migrationPath := "migrations/000001_init.up.sql"
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		logger.Error("Failed to read migration file", "path", migrationPath, "error", err)
		os.Exit(1)
	}

	_, err = dbPool.Exec(ctx, string(migrationSQL))
	if err != nil {
		logger.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Migrations applied successfully")

	secret := os.Getenv("JWT_SECRET")
	ttl, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		logger.Error("Invalid JWT_TTL format", "error", err)
		os.Exit(1)
	}

	tokenManager, err := jwt.NewTokenManager(secret, ttl)
	if err != nil {
		logger.Error("Failed to initialize token manager", "error", err)
		os.Exit(1)
	}

	repo := postgres.NewUserRepository(dbPool)
	authUseCase := usecase.NewAuthInteractor(repo, tokenManager)
	handler := delivery.NewHandler(authUseCase)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", handler.Register)
	mux.HandleFunc("/auth/login", handler.Login)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	wrappedMux := delivery.LoggingMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      wrappedMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("Server is running", "port", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}