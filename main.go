package main

import (
	"fmt"
	"github.com/fichca/image-loader/internal/config"
	"github.com/fichca/image-loader/internal/repository"
	"github.com/fichca/image-loader/internal/server"
	"github.com/fichca/image-loader/internal/service"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {

	logger := logrus.New()
	cfg := getConfig(logger)

	dbConnection := getConnectionToDB(cfg.DB, logger)

	userRepo := getUserRepo(dbConnection, cfg.DB, logger)
	userService := service.NewUserService(userRepo)

	srv := server.NewUserHandler(cfg.App.Port, logger, userService)

	srv.RegisterRoutes()
	srv.StartServer()
}

func getUserRepo(db *sqlx.DB, cfg *config.DB, logger *logrus.Logger) *repository.UserRepo {
	userRepo := repository.NewUserRepo(db, cfg)
	err := userRepo.RunMigrations()
	if err != nil {
		logger.Warning(err)
	}
	return userRepo
}

func getConfig(logger *logrus.Logger) *config.Config {
	cfg := config.Config{}
	err := cfg.Process()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(cfg.DB.Driver)
	return &cfg
}

func getConnectionToDB(cfg *config.DB, logger *logrus.Logger) *sqlx.DB {
	db, err := sqlx.Connect(cfg.Driver, fmt.Sprintf("user=%s dbname=%s sslmode=%s password=%s", cfg.User,
		cfg.Name, cfg.SSLMode, cfg.Password))
	if err != nil {
		logger.Fatal(err)
	}
	return db
}
