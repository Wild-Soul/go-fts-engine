package main

import (
	"flag"
	"log"
	"os"
	"time"

	utils "github.com/Wild-Soul/go-fts-engine/utils"
)

func main() {
	var dumpPath, queryString string
	// Getting the data file from S3 would be more interesting.
	flag.StringVar(&dumpPath, "p", "./data/enwiki-latest-abstract1.xml.gz", "Path of the wiki data")
	flag.StringVar(&queryString, "q", "Small wild cat", "search query")
	flag.Parse()
	log.Println("Flags::", dumpPath, queryString)
	log.Println("Full text search is in progress")

	// 1. Load the documents
	start := time.Now()
	docs, err := utils.LoadDocuments(dumpPath)
	if err != nil {
		log.Fatalf("Error while reading documents: %v\n", err)
		os.Exit(7) // siuuu
	}
	log.Printf("Loaded %d docs in %v", len(docs), time.Since(start))

	// 2. Index the documents.
	searchIdx := make(utils.Index)
	searchIdx.Add(docs)
	log.Printf("Indexed %d docs in %v\n", len(docs), time.Since(start))
	start = time.Now()

	// 3. Search
	matchedIds := searchIdx.Search(queryString)
	log.Printf("Search found %d ids in %v", len(matchedIds), time.Since(start))
	for _, id := range matchedIds {
		log.Printf("id: %d\ttext: %v", id, docs[id].Text)
	}

}
