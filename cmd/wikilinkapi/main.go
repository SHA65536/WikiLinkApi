package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/sha65536/wikilinkapi/wikilinkapi"
	"github.com/urfave/cli/v2"
)

func main() {
	var config_path string
	// Default configuration
	var (
		// DB_PATH
		db_path string = "/var/wikilinkapi/heb_articleslinks.db"
		// API_PORT
		port string = "2048"
		// LOG_LEVEL
		loglevel string = "info"
		// LOG_PATH
		logpath string = "/var/wikilinkapi/wikilinkapi.log"
	)

	app := &cli.App{
		Name:  "wikilinkapi",
		Usage: "Serves the WikiLink api",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config_path",
				Aliases:     []string{"c"},
				Value:       "/var/wikilinkapi/.env",
				Usage:       "Path to config file",
				Destination: &config_path,
			},
		},
		Action: func(ctx *cli.Context) error {
			// Loading environment variables
			if err := godotenv.Load(config_path); err != nil {
				return err
			}
			if val, ok := os.LookupEnv("DB_PATH"); ok {
				db_path = val
			}
			if val, ok := os.LookupEnv("API_PORT"); ok {
				port = val
			}
			if val, ok := os.LookupEnv("LOG_LEVEL"); ok {
				loglevel = val
			}
			if val, ok := os.LookupEnv("LOG_PATH"); ok {
				logpath = val
			}
			// Making logger
			var level, err = zerolog.ParseLevel(loglevel)
			if err != nil {
				return err
			}
			logf, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			// Creating API
			api, err := wikilinkapi.MakeApiHandler(db_path, level, logf)
			if err != nil {
				return err
			}
			// Serving
			err = api.Serve(":" + port)
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
