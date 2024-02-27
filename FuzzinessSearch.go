package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

/*
Fuzziness is a feature of Elasticsearch that allows for approximate matches to be found.
This can be useful for dealing with typos, misspellings, and other errors in search terms.
*/

type FuzzyQuerySearch struct {
	Index     string `json:"index"`
	Field     string `json:"field"`
	Value     string `json:"value"`
	Fuzzniess int    `json:"fuzzniess"`
}

func fuzzysearchData(c *gin.Context) {
	fmt.Println("################# ", esc)
	var fuzzyQuery FuzzyQuerySearch

	if err := c.BindJSON(&fuzzyQuery); err != nil {
		fmt.Print("Error with input", fuzzyQuery)
		return
	}
	//create fuzzy query
	query := ""
	/*query := `{` +
	`"query": {` +
	`"fuzzy": {` +
	`"` + fuzzyQuery.Field + `": {` +
	`"value": "` + fuzzyQuery.Field + `",` +
	`"fuzziness": ` + fuzzyQuery.Fuzzniess +
	`}}}` +
	`}`*/
	fmt.Printf(query)
	req := esapi.SearchRequest{
		Index: []string{fuzzyQuery.Index},
		Body:  strings.NewReader(query),
	}
	res, err := req.Do(context.Background(), esc)
	if err != nil {
		log.Fatal("Error searching documents : s ", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("No data in the body")
	}
	//jsonString, _ := json.Marshal(res)
	var results map[string]interface{}
	json.Unmarshal(body, &results)
	c.IndentedJSON(http.StatusCreated, results)
}
