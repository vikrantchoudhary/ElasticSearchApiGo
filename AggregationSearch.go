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
	"github.com/mottaquikarim/esquerydsl"
)

// BoolQuery Elastic bool query
type BoolQuery struct {
	Bool BoolQueryParams `json:"bool"`
}

// BoolQueryParams params for an Elastic bool query
type BoolQueryParams struct {
	Must               interface{} `json:"must,omitempty"`
	Should             interface{} `json:"should,omitempty"`
	Filter             interface{} `json:"filter,omitempty"`
	MinimumShouldMatch int         `json:"minimum_should_match,omitempty"`
}

/*
// The JSON representation of the BoolQueryParams struct is:
// {
//   "must": [{"term": {"name": "John Doe"}}],
//   "should": [{"term": {"age": 30}}],
//   "minimum_should_match": 1
// }
*/
type SearchAggQuery struct {
	Index     string      `json:"index"`
	Field     string      `json:"field"`
	Value     string      `json:"value"`
	QueryType string      `json:"type"`
	Sort      interface{} `json:"sort,omitempty"`
	Must      interface{} `json:"must,omitempty"`
	Should    interface{} `json:"should,omitempty"`
}

func AggrestionSearch(c *gin.Context) {
	fmt.Println("################# ", esc)
	var searchAggQuery SearchAggQuery
	if err := c.BindJSON(&searchAggQuery); err != nil {
		fmt.Print("Error with input", searchAggQuery)
		return
	}
	/*query := BoolQueryParams{}
	query.Must = searchAggQuery.Must
	query.Should = searchAggQuery.Should
	query.MinimumShouldMatch = 1

	if err != nil {
		fmt.Println("Error with input", query)
	}*/

	query, _ := json.Marshal(esquerydsl.QueryDoc{
		Index: searchAggQuery.Index,
		Sort:  []map[string]string{map[string]string{"_score": "asc"}},
		And: []esquerydsl.QueryItem{
			esquerydsl.QueryItem{
				Field: searchAggQuery.Field,
				Value: searchAggQuery.Value,
				Type:  esquerydsl.Match,
			},
		},
	})

	/*GetQueryBlock(esquerydsl.QueryDoc{
		Index: searchAggQuery.Index,
		Sort:  []map[string]string{map[string]string{"id": "asc"}},
		And: []esquerydsl.QueryItem{
			esquerydsl.QueryItem{
				Field: "some_index_id",
				Value: "some-long-key-id-value",
				Type:  "match",
			},
		},
	})*/
	fmt.Println(string(query))
	req := esapi.SearchRequest{
		Body: strings.NewReader(string(query)),
	}
	//fmt.Println(query)
	//c.IndentedJSON(http.StatusCreated, string(body))

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
