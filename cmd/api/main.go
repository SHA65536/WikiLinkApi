package main

import (
	"github.com/sha65536/wikilink/linkapi"
)

func main() {
	// Creating api
	api, err := linkapi.MakeApiHandler("bolt.db")
	if err != nil {
		panic(err)
	}

	// Starting api
	err = api.Serve(":2048")
	if err != nil {
		panic(err)
	}
}
