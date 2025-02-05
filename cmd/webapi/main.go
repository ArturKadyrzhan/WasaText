package main

import (
	"WasaText/service/api"
	"WasaText/service/database/gorm"
	"WasaText/service/repositories"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

func run() error {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	_, err := LoadConfiguration()
	db, err := gorm.NewGormSqliteDB()

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repository := repositories.NewRepository(db)
	handler := api.NewHandler(repository)

	srv := new(Server)
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
	return nil
}
