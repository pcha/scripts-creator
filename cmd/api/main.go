package main

import (
	"log"

	"scripts-creator/cmd/api/scripts"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := buildConfig()
	h := Handler{
		creator: scripts.GetCreator(cfg.GenDir),
	}

	r := setupRouter(h)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func setupRouter(h Handler) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/scripts", h.createScript)
	return r
}
