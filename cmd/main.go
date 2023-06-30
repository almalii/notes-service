package main

import (
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"log"
	"notes-rew/app"
	"notes-rew/internal/config"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	gob.Register(uuid.UUID{})

	if err := config.Loader(); err != nil {
		log.Fatalln(err)
	}

	newApp := app.NewApp()
	if err := newApp.Start(config.ServerPort()); err != nil {
		log.Fatalln(err)
	}

}
