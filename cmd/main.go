package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/app"
	"notes-rew/internal/config"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	cfg := config.InitConfig()
	ctx := context.Background()

	newApp := app.NewApp(ctx, cfg)
	if err := newApp.Start(); err != nil {
		logrus.Fatalf("Failed to run app: %+v", err)
	}

}
