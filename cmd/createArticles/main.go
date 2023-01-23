package main

import (
	"log"

	"github.com/sha65536/linkapibe/linkapibe"
)

func main() {
	// Creating db handler
	db, err := linkapibe.MakeDbHandler("bolt.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.CreateBuckets()
	if err != nil {
		panic(err)
	}

	// Creating scrape handler
	var count uint32
	sc := linkapibe.MakeScrapeHandler(db)
	err = sc.GetAllArticles("he", func(s []string) {
		for i := range s {
			if err := sc.DB.CreateArticle(s[i], count); err != nil {
				log.Fatal(err)
			}
			count++
		}
		log.Println(count)
	})
	if err != nil {
		panic(err)
	}
}
