package main

import (
	"github.com/sirupsen/logrus"
	"notes-rew/internal/app"
	"notes-rew/internal/config"
)

func main() {

	cfg := config.InitConfig()

	if err := app.NewApp().Start(cfg.HTTPServer.Port); err != nil {
		logrus.Fatal(err)
	}

}
