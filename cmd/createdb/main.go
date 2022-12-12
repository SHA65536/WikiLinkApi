package main

import "github.com/sha65536/wikilink/linkapi"

func main() {
	db, err := linkapi.MakeDbHandler("bolt.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}
