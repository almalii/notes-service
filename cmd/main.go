package main

import (
	"context"
	"encoding/gob"
	"notes-rew/internal/app"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/config"
)

// @title Notes-Service API
// @version 1.0
// @description This is a sample notes-rew server.

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apiKey JWTAuth
// @in header
// @name Authorization
func main() {
	gob.Register(uuid.UUID{})
	logrus.SetFormatter(&logrus.JSONFormatter{})

	cfg := config.InitConfig("local")
	ctx := context.Background()

	newApp := app.NewApp(ctx, *cfg)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := newApp.StartGRPC(); err != nil {
			logrus.Fatalf("Failed to start GRPC server: %+v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := newApp.Start(ctx); err != nil {
			logrus.Fatalf("Failed to start HTTP server: %+v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := newApp.StartGateway(ctx); err != nil {
			logrus.Fatalf("Failed to start GRPC-Gateway server: %+v", err)
		}
	}()

	wg.Wait()

}
