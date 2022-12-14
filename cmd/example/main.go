package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	// Creating search handler
	sc := linkapi.MakeSearchHandler(db)

	// Getting all articles map
	arts, err := db.GetAllArticles()
	if err != nil {
		panic(err)
	}

	// Showing names to id
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		if text == "\n" {
			break
		}
		if val, ok := arts[text[:len(text)-1]]; ok {
			fmt.Printf("article %s is id %d\n", text[:len(text)-1], val)
		}
	}
	// Reading two ids
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter two article ids separated by a space:")
	text, _ := reader.ReadString('\n')
	ids := strings.Split(text, " ")
	id1, _ := strconv.ParseUint(ids[0], 10, 32)
	id2, _ := strconv.ParseUint(ids[1][:len(ids[1])-1], 10, 32)
	fmt.Println(id1, id2)
	// Searching
	res, err := sc.ShortestPath(uint32(id1), uint32(id2), func(i int) {})
	if err != nil {
		panic(err)
	}
	// Converting ids to names
	for i := range res {
		v, err := db.GetName(res[i])
		if err != nil {
			panic(err)
		}
		fmt.Println(v)
	}
}
