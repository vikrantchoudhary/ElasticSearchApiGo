package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gin-gonic/gin"
)

func LoadDataFromURL(c *gin.Context) {
	var url string
	if err := c.BindJSON(&url); err != nil {
		return
	}
	var spacecrafts []map[string]interface{}
	pageNumber := 3
	for {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println("No data fetched from the url " + url)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println("No data in the body" + url)
		}
		defer res.Body.Close()
		var results map[string]interface{}
		json.Unmarshal(body, &results)
		page := results["page"].(map[string]interface{})
		totalPages := int(page["totalPages"].(float64))
		crafts := results["spacecrafts"].([]interface{})
		for _, craftInterface := range crafts {
			craft := craftInterface.(map[string]interface{})
			spacecrafts = append(spacecrafts, craft)
		}
		pageNumber++
		if pageNumber >= totalPages {
			break
		}
	}
	for _, data := range spacecrafts {
		uid, _ := data["uid"].(string)
		jsonString, _ := json.Marshal(data)
		res := esapi.IndexRequest{
			Index:      "spacecraft",
			DocumentID: uid,
			Body:       strings.NewReader(string(jsonString)),
			Refresh:    "true",
		}
		res.Do(context.Background(), esc)
	}
	c.IndentedJSON(http.StatusCreated, "Indices created")

}
