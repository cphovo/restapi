package main

import (
	"github.com/cphovo/restapi/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// GenerateImage(600, 600, "#FF5733")
	r := gin.Default()
	r.GET("/", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Everything is ok!"}) })
	routes.SetupRoutes(r)
	r.Run(":8000")
}
