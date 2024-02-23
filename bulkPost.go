package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

const (
	elasticIndexName = "documents"
	elasticTypeName  = "document"
)

type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

type DocumentRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type DocumentResponse struct {
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

func createDocumentsEndpoint(c *gin.Context) {
	var docs []DocumentRequest
	if err := c.BindJSON(&docs); err != nil {
		errorResponse(c, http.StatusBadRequest, "Malformed request body")
		return
	}
	fmt.Print("################# ", esc)
	// Insert documents in bulk
	bulk, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         elasticIndexName,
		Client:        esc,
		NumWorkers:    runtime.NumCPU(),
		FlushBytes:    int(5e+6),
		FlushInterval: 30 * time.Second,
	})
	if err != nil {
		log.Fatal("Error creating the indexer: ", err)
	}

	for _, d := range docs {
		doc := Document{
			ID:        shortid.MustGenerate(),
			Title:     d.Title,
			CreatedAt: time.Now().UTC(),
			Content:   d.Content,
		}
		data, err := json.Marshal(d)
		if err != nil {
			log.Fatalf("Cannot encode article %s: %s", doc.ID, err)
		}
		//bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
		err = bulk.Add(
			context.Background(),

			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: doc.ID,
				Body:       bytes.NewReader(data),
				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
					//fmt.Println("processed ", doc.ID)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatal("Unexpected error: ", err)
		}
	}
	if err != nil {
		log.Fatal("Unexpected error: ", err)
	}
	c.Status(http.StatusOK)
}
