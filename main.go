package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
	_ "github.com/yousefggg/auth-service/docs"
	"github.com/yousefggg/auth-service/internal/delivery"
	"github.com/yousefggg/auth-service/internal/repository/postgres"
	"github.com/yousefggg/auth-service/internal/usecase"
	"github.com/yousefggg/common-lib/pkg/jwt"
	"github.com/yousefggg/common-lib/pkg/logger"
)

// @title           Mountain Tour Auth Service API
// @version         1.0
// @description     Сервис аутентификации и управления пользователями.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Yousef Mulaev
// @contact.email  mulaev2006@gmail.com

// @host      localhost:8081
// @BasePath  /
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

	dsn := os.Getenv("DATABASE_URL")
	runMigrations(dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, dsn)
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

func runMigrations(dsn string) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Error("Could not create migrate instance", "error", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Migrations applied successfully")
	}
}