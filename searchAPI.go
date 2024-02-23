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

type QueryString struct {
	Index string `json:"index"`
	Field string `json:"field"`
	Value string `json:"value"`
}

func searchData(c *gin.Context) {
	fmt.Println("################# ", esc)
	var queryString QueryString

	if err := c.BindJSON(&queryString); err != nil {
		fmt.Print("Error with input", queryString)
		return
	}
	//fmt.Println(queryString.Field + queryString.Index + queryString.Value)
	query := `{"query": {"match" : {"` + queryString.Field + `":"` + queryString.Value + `"}}}`

	req := esapi.SearchRequest{
		Index: []string{queryString.Index},
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
