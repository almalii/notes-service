package migrations

import (
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/config"
)

func UpMigrations(config *config.MigrationsConfig) error {
	gooseDB, err := goose.OpenDBWithDriver(config.Driver, config.ConnString)
	if err != nil {
		logrus.Error("error opening db connection on migrations")
		return err
	}

	defer gooseDB.Close()

	err = goose.SetDialect(config.Driver)
	if err != nil {
		logrus.Error("setting dialect error on migrations")
		return err
	}

	err = goose.Up(gooseDB, config.MigrationsDir)
	if err != nil {
		logrus.Error("running migrations error on up")
		return err
	}

	return nil
}

func DownMigrations(config *config.MigrationsConfig) error {
	gooseDB, err := goose.OpenDBWithDriver(config.Driver, config.ConnString)
	if err != nil {
		logrus.Error("error opening db connection on migrations")
		return err
	}

	defer gooseDB.Close()

	err = goose.SetDialect(config.Driver)
	if err != nil {
		logrus.Error("setting dialect error on migrations")
		return err
	}

	err = goose.Down(gooseDB, config.MigrationsDir)
	if err != nil {
		logrus.Error("running migrations error on down")
		return err
	}

	return nil
}
