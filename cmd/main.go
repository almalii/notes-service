package main

import (
	"context"
	"encoding/gob"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/app/grpc_app"
	"notes-rew/internal/app/rest_app"
	"notes-rew/internal/config"
	"sync"
)

///@title Notes-rew API
///@version 1.0
///@description This is a sample notes-rew server.

///@host localhost:8080

func main() {
	gob.Register(uuid.UUID{})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	//logrus.SetLevel(logrus.DebugLevel)

	cfg := config.InitConfig()
	ctx := context.Background()

	newApp := rest_app.NewApp(ctx, cfg)
	newAppGRPC := grpc_app.NewAppGRPC(ctx, cfg)

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		if err := newAppGRPC.StartGRPC(); err != nil {
			logrus.Fatalf("Не удалось запустить GRPC-сервер: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newApp.Start(); err != nil {
			logrus.Fatalf("Не удалось запустить HTTP сервер: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newAppGRPC.StartGateway(); err != nil {
			logrus.Fatalf("Не удалось запустить GRPC-Gateway сервер: %+v", err)
		}
	}()

	wg.Wait()

}
