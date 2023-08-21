package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"path"
	"runtime"
	"sync"

	"notes-rew/internal/app"

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

	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: false,
		FullTimestamp: true,
	})

	cfg := config.InitConfig()

	logrus.Info("cfg.SaltHash", cfg.SaltHash)

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
