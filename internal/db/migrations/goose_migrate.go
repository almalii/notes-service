package migrations

import (
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/config"
)

func RunMigrations(config *config.MigrationsConfig, command string) error {
	gooseDB, err := goose.OpenDBWithDriver(config.Driver, config.ConnString)
	if err != nil {
		return err
	}

	defer gooseDB.Close()

	err = goose.SetDialect(config.Driver)
	if err != nil {
		return err
	}

	switch command {
	case "up":
		err = goose.Up(gooseDB, config.MigrationsDir)
	case "down":
		err = goose.Down(gooseDB, config.MigrationsDir)
	default:
		logrus.Errorf("Unknown command. Usage: up or down")
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
