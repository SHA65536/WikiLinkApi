package main

import (
	"log"
	"os"

	"github.com/sha65536/linkapibe/linkapibe"
	"github.com/urfave/cli/v2"
)

func main() {
	var port string
	var db_path string

	app := &cli.App{
		Name:        "linkapi",
		Description: "Serves the link api",
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
		},
		Action: func(ctx *cli.Context) error {
			api, err := linkapibe.MakeApiHandler(db_path)
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
