package main

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func initializeElasticSearchConfig() {
	//dsn := "sql6684931:NQFt2Xq66u@tcp(sql6.freemysqlhosting.net:3306)/sql6684931?charset=utf8mb4&parseTime=True&loc=Local"
	err := godotenv.Load("setup.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	esUsername = os.Getenv("ES_USERNAME")
	esPassword = os.Getenv("ES_PASSWORD")
	esCertPath = os.Getenv("HTTP_CERT")
	esURL = os.Getenv("ES_URL")

}

/*
	func initializeElasticSearchConfig() {
		cert, _ := os.ReadFile("/Users/vikrant/http_ca.crt")
		cfg := elasticsearch.Config{
			Addresses: []string{
				"https://localhost:9200",
			},
			CACert:   cert,
			Username: "elastic",
			Password: "Ilyw_HKZCw+PSFDajCtz",
		}
		esc, err := elasticsearch.NewClient(cfg)
		if err != nil {
			fmt.Println(err.Error())
			log.Fatal("Cannot connect to Elastic Search with the given credential")
		}
		if esc != nil {
			fmt.Print(esc)
		}
	}
*/
func Init() {
	cert, _ := os.ReadFile(esCertPath)
	cfg := elasticsearch.Config{
		Addresses: []string{
			esURL,
		},
		CACert:   cert,
		Username: esUsername,
		Password: esPassword,
	}
	esc, err = elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Cannot connect to Elastic Search with the given credential")
	}
}

func errorResponse(c *gin.Context, code int, err string) {
	c.JSON(code, gin.H{
		"error": err,
	})
}
