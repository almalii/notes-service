package main

import (
	"log"
	"notes-rew/internal/app"
	"notes-rew/internal/config"
)

func main() {

	if err := config.Loader(); err != nil {
		log.Fatalln(err)
	}

	newApp := app.NewApp()
	if err := newApp.Start(config.ServerPort()); err != nil {
		log.Fatalln(err)
	}

}
