package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fichca/image-loader/internal/config"
	"github.com/fichca/image-loader/internal/filestore"
	"github.com/fichca/image-loader/internal/middleware"
	"github.com/fichca/image-loader/internal/repository"
	"github.com/fichca/image-loader/internal/server"
	"github.com/fichca/image-loader/internal/service"
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logrus.New()
	router := chi.NewRouter()

	router.Use(middleware.Logger(logger))

	cfg := initConfig(logger)

	userRepo, imageRepo, fileStorage := initRepositories(logger, cfg)

	authService := service.NewAuthService(userRepo, cfg.JWTKeyword)
	fileService := service.NewFileService(fileStorage, imageRepo)
	userService := service.NewUserService(userRepo, fileService)

	authMiddleware := middleware.Auth(authService, cfg.JWTKeyword, logger)

	userHandler := server.NewUserHandler(logger, userService, router, authMiddleware)
	userHandler.RegisterUserRoutes()

	fileHandler := server.NewFileHandler(logger, fileService, router, authMiddleware)
	fileHandler.RegisterFileRoutes()

	authHandler := server.NewAuthHandler(logger, authService, router)
	authHandler.RegisterAuthRoutes()

	startServer(cfg.App.Port, router, logger)
}

func initRepositories(logger *logrus.Logger, cfg *config.Config) (*repository.UserRepo, *repository.ImageRepo, *filestore.Minio) {
	dbConnection := initDBConnection(cfg.DB, logger)
	minioConnection := initMinioConnection(logger, cfg.Minio)
	userRepo := repository.NewUserRepo(dbConnection)
	imageRepo := repository.NewImageRepo(dbConnection)
	fileStorage := filestore.NewMinio(minioConnection, cfg.Minio.Bucket)
	err := RunMigrations(dbConnection.DB, cfg)
	if err != nil {
		logger.Warning(err)
	}
	return userRepo, imageRepo, fileStorage
}

func startServer(listenURI string, r chi.Router, logger *logrus.Logger) {
	srv := http.Server{
		Addr:    listenURI,
		Handler: r,
	}

	logger.Info(fmt.Sprintf("server is running on port %v!", listenURI))
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}

func initMinioConnection(logger *logrus.Logger, cfg *config.Minio) *minio.Client {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.KeyID, cfg.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		logger.Fatal(err)
	}

	ok, err := minioClient.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		logger.Fatal(err)
	}

	if !ok {
		err = minioClient.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Fatal(err)
		}
	}
	return minioClient
}

func initConfig(logger *logrus.Logger) *config.Config {
	cfg := config.Config{}
	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(cfg.DB.Driver)
	return &cfg
}

func initDBConnection(cfg *config.DB, logger *logrus.Logger) *sqlx.DB {
	db, err := sqlx.Connect(cfg.Driver, fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s", cfg.User,
		cfg.Name, cfg.SSLMode, cfg.Password))
	if err != nil {
		logger.Fatal(err)
	}
	return db
}

func RunMigrations(dbConnection *sql.DB, cfg *config.Config) error {
	driver, err := postgres.WithInstance(dbConnection, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get migration tool driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		cfg.DB.Driver, driver)
	if err != nil {
		return fmt.Errorf("failed to connect migration tool: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
