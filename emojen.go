package main

import "fmt"
import "net/http"
import "github.com/gin-gonic/gin"

func setup_router(r *gin.Engine) {
	r.GET("/r/:link", func(c *gin.Context) {
		link := c.Param("link")
		fmt.Printf("link is %s\n", link)

		// TODO: lookup link, if found 301, if not 404
		c.Redirect(http.StatusMovedPermanently, "https://google.com/")
	})

	r.POST("/r", func(c *gin.Context) {
		link := c.PostForm("link")
		fmt.Printf("link is %s\n", link)

		// TODO: calculate short link
		// TODO: insert into redis
		short := "/r/🐯🐁🙃"

		c.JSON(200, gin.H{
			"short": short,
		})
	})
}

func main() {
	r := gin.Default()
	setup_router(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
