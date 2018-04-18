package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()

	f, _ := os.Create("/var/log/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// ws test
	r.GET("/ws", func(c *gin.Context) {
		c.HTML(http.StatusOK, "ws.html", gin.H{
			"title": "Main website",
		})
	})

	r.StaticFS("/static", http.Dir("static"))

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":5701")
}
