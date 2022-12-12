package server

import "github.com/gin-gonic/gin"

func healthHandler(c *gin.Context) {
	c.JSON(200, "ok")
}
