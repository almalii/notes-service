package main

import (
	"bookmarks/app"
	"bookmarks/internal/config"
	"encoding/gob"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
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
