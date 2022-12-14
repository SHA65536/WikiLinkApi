package main

import (
	"github.com/schollz/progressbar/v3"

	"github.com/sha65536/wikilink/linkapi"
)

func main() {
	// Creating db handler
	db, err := linkapi.MakeDbHandler("bolt.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.CreateBuckets()
	if err != nil {
		panic(err)
	}

	// Getting all articles map
	arts, err := db.GetAllArticles()
	if err != nil {
		panic(err)
	}

	// Creating scrape handler
	sc := linkapi.MakeScrapeHandler(db)
	bar := progressbar.Default(int64(len(arts)))

	// Looping over all the articles
	for k, v := range arts {
		bar.Add(1)
		// Getting all links from each
		err = sc.GetAllLinks("he", k, func(s []string) {
			var res = make([]uint32, 0, len(s))
			for i := range s {
				if id, ok := arts[s[i]]; ok {
					res = append(res, id)
				}
			}
			// Adding to DB
			if err := db.AddLinks(v, res); err != nil {
				panic(err)
			}
		})
	}
	if err != nil {
		panic(err)
	}
}
