package migrations

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/config"
)

func UpMigrations(c *config.DbConfig) error {
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s", c.Host, c.Port, c.DBName, c.SSLMode)
	gooseDB, err := goose.OpenDBWithDriver(c.Driver, connStr)
	if err != nil {
		logrus.Error("error opening db connection on migrations")
		return err
	}

	defer gooseDB.Close()

	err = goose.SetDialect(c.Driver)
	if err != nil {
		logrus.Error("setting dialect error on migrations")
		return err
	}

	err = goose.Up(gooseDB, "../internal/db/migrations/")
	if err != nil {
		logrus.Error("running migrations error on up")
		return err
	}

	return nil
}
