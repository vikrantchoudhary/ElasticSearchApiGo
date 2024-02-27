package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func DSLQuerySearch(c *gin.Context) {
	results := "result processed"
	c.IndentedJSON(http.StatusCreated, results)
}
