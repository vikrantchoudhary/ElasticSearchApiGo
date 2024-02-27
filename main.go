package main

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

var (
	esc             *elasticsearch.Client
	countSuccessful uint64
	err             error
)
var (
	esUsername string
	esPassword string
	esURL      string
	esCertPath string
)

func main() {
	initializeElasticSearchConfig()
	Init()
	r := gin.Default()
	r.POST("/documents", createDocumentsEndpoint)
	r.POST("/url", LoadDataFromURL)
	r.POST("/search", searchData)
	if err = r.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
