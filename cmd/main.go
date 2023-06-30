package main

import (
	"bookmarks/app"
	"bookmarks/internal/config"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"log"
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
