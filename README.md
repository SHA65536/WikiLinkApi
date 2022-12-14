# WikiLink
This is a repository for scraping, storing, and searching links between wikipedia articles.

# Scraping
`scrape.go` contains methods for getting all articles from a regional wikipedia, and getting all links from a list of articles.

# Storing
`database.go` contains a key value database handler using [bbolt](go.etcd.io/bbolt) as it's underlying database. It allows efficient saving and querying of articles and links.

# Searching
`search.go` contains a BFS implementation to find the shortest path between two articles in the database.

# Scripts
The `cmd` directory contains some useful scripts for scraping articles and links, example for shortest path and an api serving script.