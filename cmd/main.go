package main

import (
	"WasaText/cmd/database/gorm"
	"WasaText/cmd/webapi"
	"WasaText/internal/api"
	"WasaText/internal/repositories"
	"WasaText/internal/service"
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("error loading .env file: %v", err.Error())
	}
	db, err := gorm.NewGormSqliteDB()

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repository := repositories.NewRepository(db)
	services := service.NewService(repository)
	handler := api.NewHandler(services)

	srv := new(webapi.Server)
	go func() {
		if err := srv.Run(os.Getenv("SERVER_PORT"), handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Project is Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Project is Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}
