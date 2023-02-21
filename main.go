package main

import (
	"github.com/sirupsen/logrus"
	"image-loader/internal/repository"
	"image-loader/internal/server"
	"image-loader/internal/services"
)

func main() {
	userRepo := repository.NewUserRepo("test.json")
	userService := services.NewUserService(userRepo)

	logger := logrus.New()

	srv := server.NewUserHandler(":8080", logger, userService)
	srv.RegisterRoutes()

	srv.StartServer()
}
