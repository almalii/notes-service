package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/app"
	"notes-rew/internal/config"
	"sync"
)

func main() {

	logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetLevel(logrus.DebugLevel)

	cfg := config.InitConfig()
	ctx := context.Background()

	newApp := app.NewApp(ctx, cfg)
	newAppGRPC := app.NewAppGRPC(ctx, cfg)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := newAppGRPC.StartGRPC(); err != nil {
			logrus.Fatalf("Не удалось запустить GRPC-сервер: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newApp.StartHTTP(); err != nil {
			logrus.Fatalf("Не удалось запустить приложение: %+v", err)
		}
	}()

	wg.Wait()

}
