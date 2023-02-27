package main

import (
	"fmt"
	"github.com/fichca/image-loader/internal/config"
	"github.com/fichca/image-loader/internal/middleware"
	"github.com/fichca/image-loader/internal/repository"
	"github.com/fichca/image-loader/internal/server"
	"github.com/fichca/image-loader/internal/service"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logger := logrus.New()
	router := chi.NewRouter()
	router.Use(middleware.Logger(logger))

	cfg := initConfig(logger)

	dbConnection := initConnectionToDB(cfg.DB, logger)

	userRepo := initUserRepo(dbConnection, cfg.DB, logger)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTKeyword)

	userHandler := server.NewUserHandler(logger, userService, authService, router, cfg.JWTKeyword)
	userHandler.RegisterUserRoutes()

	authHandler := server.NewAuthHandler(logger, authService, router)
	authHandler.RegisterAuthRoutes()

	startServer(cfg.App.Port, router, logger)
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

func initUserRepo(db *sqlx.DB, cfg *config.DB, logger *logrus.Logger) *repository.UserRepo {
	userRepo := repository.NewUserRepo(db, cfg)
	err := userRepo.RunMigrations()
	if err != nil {
		logger.Warning(err)
	}
	return userRepo
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

func initConnectionToDB(cfg *config.DB, logger *logrus.Logger) *sqlx.DB {
	db, err := sqlx.Connect(cfg.Driver, fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s", cfg.User,
		cfg.Name, cfg.SSLMode, cfg.Password))
	if err != nil {
		logger.Fatal(err)
	}
	return db
}
