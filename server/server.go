package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func ServeGin(authCodeChan chan string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(gin.Recovery())

	r.GET("/auth", func(ctx *gin.Context) {
		code := ctx.Query("code")
		authCodeChan <- code
		ctx.JSON(200, gin.H{
			"status": "Access Code recieved. You can close this window.",
		})
	})

	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
