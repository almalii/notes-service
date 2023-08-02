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
			logrus.Fatalf("Failed to start GRPC server: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newApp.Start(); err != nil {
			logrus.Fatalf("Failed to start HTTP server: %+v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := newAppGRPC.StartGateway(); err != nil {
			logrus.Fatalf("Failed to start GRPC-Gateway server: %+v", err)
		}
	}()

	wg.Wait()

}
