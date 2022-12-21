# WikiLink
This is a repository for scraping, storing, and searching links between wikipedia articles.

<p>
   <img src="https://raw.githubusercontent.com/SHA65536/WikipediaProject/main/.github/wiki.gif" width="100%" alt="WikiLink">
</p>

# Scraping
`scrape.go` contains methods for getting all articles from a regional wikipedia, and getting all links from a list of articles.

# Storing
`database.go` contains a key value database handler using [bbolt](go.etcd.io/bbolt) as it's underlying database. It allows efficient saving and querying of articles and links.

# Searching
`search.go` contains a BFS implementation to find the shortest path between two articles in the database.

# Scripts
The `cmd` directory contains some useful scripts for scraping articles and links, example for shortest path and an api serving script.

# Api
To run the api run `go run ./cmd/linkapi`:
```
NAME:
   linkapi - A new cli application

USAGE:
   linkapi [global options] command [command options] [arguments...]

DESCRIPTION:
   Serves the link api

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value  Port to listen to (default: "2048")
   --db value, -d value    Path to the database (default: "bolt.db")
   --help, -h              show help (default: false)
``` 