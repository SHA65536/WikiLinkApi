# WikiLinkAPI
This is a repository for scraping, storing, and searching links between wikipedia articles.

<p>
   <img src="https://raw.githubusercontent.com/SHA65536/WikipediaProject/main/.github/wiki.gif" width="100%" alt="WikiLink">
</p>

## Scraping
`scrape.go` contains methods for getting all articles from a regional wikipedia, and getting all links from a list of articles.

## Storing
`database.go` contains a key value database handler using [bbolt](go.etcd.io/bbolt) as it's underlying database. It allows efficient saving and querying of articles and links.

## Searching
`search.go` contains a BFS implementation to find the shortest path between two articles in the database.

## Scripts
The `cmd` directory contains some useful scripts for scraping articles and links, example for shortest path and an api serving script.

## Api
To run the api run `go run ./cmd/wikilinkapi`:
```
NAME:
   wikilinkapi - Serves the WikiLink api

USAGE:
   wikilinkapi [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config_path value, -c value  Path to config file (default: "/var/wikilinkapi/.env")
   --help, -h                     show help (default: false)", "debug", "info", "warn", "error", "fatal", "panic") (default: "info")
   --help, -h              show help (default: false)
``` 

## Configuration
The .env configuration file should have the following format:
```
DB_PATH="heb_articleslinks.db"
API_PORT="2048"
LOG_LEVEL="info"
LOG_PATH="wikilinkapi.log"
```

## Database download
You can download the pre-scraped articles and links database for the hebrew wikipedia using [this link!](https://mega.nz/file/JooCEaoA#cAuECOOFXBF8oTB6410yJxqy5X4c5eL_3A_Z591I8R0)