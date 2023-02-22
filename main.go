package main

import (
	"github.com/fichca/image-loader/internal/repository"
	"github.com/fichca/image-loader/internal/server"
	"github.com/fichca/image-loader/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	userRepo := repository.NewUserRepo("test.json")
	userService := services.NewUserService(userRepo)

	logger := logrus.New()

	srv := server.NewUserHandler(":8080", logger, userService)
	srv.RegisterRoutes()

	srv.StartServer()
}
