package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mottaquikarim/esquerydsl"
)

/*
* Advance query search : full Query DSL (Domain Specific Language) based on JSON to define queries
 */

type SearchQuery struct {
	QueryString  string
	SearchFields []string
	SortField    string
	SortOrder    string
}

func searchQueryDSL(c *gin.Context) {
	fmt.Println("################# ", esc)
	var searchQuery SearchQuery
	if err := c.BindJSON(&searchQuery); err != nil {
		fmt.Print("Error with input", searchQuery)
		return
	}
	//var payload := []byte
	queries := ""
	/*
			POST /craft/_search?pretty
		  {
		    "query": {
		      "bool": {
		        "must": [
		          {"match": {
		            "name": "SS"
		          }

		          }

		        ]
		      }
		    },
		    "sort": [
		      {
		        "_score": {
		          "order": "desc"
		        }
		      }
		    ]
		  }
	*/

	body, _ := json.Marshal(esquerydsl.QueryDoc{
		Index: "some_index",
		Sort:  []map[string]string{map[string]string{"id": "asc"}},
		And: []esquerydsl.QueryItem{
			esquerydsl.QueryItem{
				Field: "some_index_id",
				Value: "some-long-key-id-value",
				Type:  esquerydsl.Match,
			},
		},
	})
	fmt.Println(string(body))
	/*for _, field := range searchQuery.SearchFields {
			query := `{"index":`+ `"}`
			queries = append(queries,)
	        esQuery = esQuery.Should(elastic.NewMatchQuery(field, query.QueryString))
	    }*/
	//eqQuery := esapi.MsearchRequest()
	c.IndentedJSON(http.StatusCreated, queries)
}
