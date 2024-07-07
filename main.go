package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	utils "github.com/Wild-Soul/go-fts-engine/utils"
)

type Document struct {
	Id    string `json:"id"`
	Title string `json:"content"`
	Text  string `json:"tags"`
}

type Query struct {
	Query string `json:"query"`
}

func main() {

	// 1. Load the documents -- shoudl be moved to init()
	start := time.Now()
	docs, err := utils.LoadDocuments("./data/enwiki-latest-abstract1.xml.gz")
	if err != nil {
		log.Fatalf("Error while reading documents: %v\n", err)
		os.Exit(7) // siuuu
	}
	log.Printf("Loaded %d docs in %v", len(docs), time.Since(start))

	// 2. Index the documents.
	searchIdx := utils.NewIndex()
	searchIdx.Add(docs)
	log.Printf("Indexed %d docs in %v\n", len(docs), time.Since(start))

	// Start http server.
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		doc := v1.Group("/doc")
		{
			doc.POST("/insert", insertDoc(searchIdx))
			doc.DELETE("/delete/:id", deleteDoc(searchIdx))
			doc.POST("/query", queryDocs(searchIdx))
		}
	}

	r.Run(":8080")
}

func insertDoc(index *utils.Index) gin.HandlerFunc {
	return func(c *gin.Context) {
		var doc Document
		if err := c.BindJSON(&doc); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Insert the document into your index
		panic("Method not yet implemented")
	}
}

func deleteDoc(index *utils.Index) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Println("Deleting document:", id)
		panic("Method not yet implemented")
	}
}

func queryDocs(index *utils.Index) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query Query
		if err := c.BindJSON(&query); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		matchedIds := index.Search(query.Query)
		c.JSON(http.StatusOK, matchedIds)
	}
}
