package main

import (
	"sync"

	"github.com/schollz/progressbar/v3"

	"github.com/sha65536/wikilinkapi/wikilinkapi"
)

const NUMWORKERS = 10

func main() {
	// Creating db handler
	db, err := wikilinkapi.MakeDbHandler("bolt.db")
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
	sc := wikilinkapi.MakeScrapeHandler(db)
	bar := progressbar.Default(int64(len(arts)))

	// Creating workers
	var wg sync.WaitGroup
	var artChan = make(chan Art)
	for i := 0; i < NUMWORKERS; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			downloadWorker(sc, artChan, arts)
		}()
	}

	// Looping over all the articles
	for k, v := range arts {
		bar.Add(1)
		// Getting all links from each
		artChan <- Art{Name: k, Id: v}
	}

	// Waiting for the last ones to finish
	close(artChan)
	wg.Wait()
}

type Art struct {
	Name string
	Id   uint32
}

func downloadWorker(sc *wikilinkapi.ScrapeHandler, in chan Art, arts map[string]uint32) {
	for val := range in {
		// Getting all links
		sc.GetAllLinks("he", val.Name, func(s []string) {
			var res = make([]uint32, 0, len(s))
			for i := range s {
				if id, ok := arts[s[i]]; ok {
					res = append(res, id)
				}
			}
			// Adding to DB
			if err := sc.DB.AddLinks(val.Id, res); err != nil {
				panic(err)
			}
		})
	}
}
