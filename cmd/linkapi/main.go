package main

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/sha65536/wikilinkapi/wikilinkapi"
	"github.com/urfave/cli/v2"
)

func main() {
	var port string
	var db_path string
	var loglevel string

	app := &cli.App{
		Name:  "linkapi",
		Usage: "Serves the WikiLink api",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       "2048",
				Usage:       "Port to listen to",
				Destination: &port,
			},
			&cli.StringFlag{
				Name:        "db",
				Aliases:     []string{"d"},
				Value:       "bolt.db",
				Usage:       "Path to the database",
				Destination: &db_path,
			},
			&cli.StringFlag{
				Name:        "log",
				Aliases:     []string{"l"},
				Value:       "info",
				Usage:       `Level of log to be shown ("trace", "debug", "info", "warn", "error", "fatal", "panic")`,
				Destination: &loglevel,
			},
		},
		Action: func(ctx *cli.Context) error {
			var level, err = zerolog.ParseLevel(loglevel)
			if err != nil {
				return err
			}
			logf, err := os.OpenFile("linkapi.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			api, err := wikilinkapi.MakeApiHandler(db_path, level, logf)
			if err != nil {
				return err
			}
			err = api.Serve(":" + port)
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
