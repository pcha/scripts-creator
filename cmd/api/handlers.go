package main

import (
	"scripts-creator/cmd/api/scripts"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	creator scripts.Creator
}

func (h Handler) createScript(c *gin.Context) {
	def := new(scripts.Definition)
	err := c.ShouldBindJSON(def)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	path2file, err := h.creator.Create(*def)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(201, gin.H{
		"location": path2file,
	})
}
