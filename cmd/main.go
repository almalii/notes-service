package main

import (
	"github.com/sirupsen/logrus"
	"notes-rew/internal/app"
	"notes-rew/internal/config"
)

func main() {

	if err := app.NewApp().Start(config.ServerPort()); err != nil {
		logrus.Fatal(err)
	}

}
