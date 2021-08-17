package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		log.Println("test")
		c.JSON(http.StatusOK,gin.H{
			"message":"success",
		})
	})

	r.Run(":8081")
}