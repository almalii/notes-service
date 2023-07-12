package mongo

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"notes-rew/internal/config"
)

func ConnectionMongoDB(ctx context.Context, c *config.Config) (*mongo.Client, error) {
	connStr := fmt.Sprintf("mongodb://%s:%s", c.DB.Host, c.DB.Port)
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))

	if c == nil {
		logrus.Fatalf("failed to connect the mongo database: %v", err)
		return nil, err
	}

	if err != nil {
		logrus.Fatalf("failed to connect to the mongo database: %v", err)
		return nil, err
	}

	if err = conn.Ping(ctx, readpref.Primary()); err != nil {
		conn.Disconnect(ctx)
		logrus.Fatalf("failed to ping the mongo database: %v", err)
		return nil, err
	}

	logrus.Println("mongo database connection successfully")

	return conn, nil
}
