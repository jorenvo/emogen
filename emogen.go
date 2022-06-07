// Copyright 2019, Joren Van Onder (joren.vanonder@gmail.com)
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func randomIndex(length uint) int {
	index := math.Floor(rand.Float64() * float64(length))
	return int(index)
}

func getShortLink(redisConn redis.Conn, shortlink string) (string, error) {
	return redis.String(redisConn.Do("GET", "shortlink:"+shortlink))
}

func randomEmojiString(length uint) string {
	return fmt.Sprintf(
		"%s%s%s",
		emojis[randomIndex(length)],
		emojis[randomIndex(length)],
		emojis[randomIndex(length)],
	)
}

func getEmojis(redisPool *redis.Pool, length uint) string {
	redisConn := redisPool.Get()
	defer redisConn.Close()

	for {
		newEmojis := randomEmojiString(length)
		link, _ := getShortLink(redisConn, newEmojis)
		if link != "" {
			log.Printf("collision for %s, retrying...\n", newEmojis)
		} else {
			return newEmojis
		}
	}
}

type toShorten struct {
	Link string `json:"link" binding:"required"`
}

func setupRouter(router *gin.Engine, redisPool *redis.Pool) {
	router.GET("/:link", func(c *gin.Context) {
		redisConn := redisPool.Get()
		defer redisConn.Close()

		shortlink := c.Param("link")

		link, err := getShortLink(redisConn, shortlink)
		if err != nil {
			link = "/"
		}

		log.Printf("resolving %s to %s\n", shortlink, link)
		c.Redirect(http.StatusMovedPermanently, link)
	})

	router.POST("/", func(c *gin.Context) {
		redisConn := redisPool.Get()
		defer redisConn.Close()

		var json toShorten
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		link := json.Link
		log.Printf("shortening %s\n", link)

		shortLink, err := redis.String(redisConn.Do("GET", "link:"+link))
		if err == nil {
			c.JSON(200, gin.H{
				"link": "/" + shortLink,
			})
			return
		}

		if !strings.HasPrefix(link, "http") {
			log.Printf("link %s did not specify http\n", link)
			link = "http://" + link
		}

		shortLink = getEmojis(redisPool, uint(len(emojis)))

		_, err = redisConn.Do("SET", "shortlink:"+shortLink, link)
		if err != nil {
			log.Printf("Error while storing link %s -> %s (%s)\n", shortLink, link, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed connecting to db.",
			})
			return
		}

		_, err = redisConn.Do("SET", "link:"+link, shortLink)
		if err != nil {
			log.Printf("Error while storing link %s -> %s (%s)\n", link, shortLink, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed connecting to db.",
			})
			return
		}

		c.JSON(200, gin.H{
			"link": "/" + shortLink,
		})
	})
}

func setup() {
	rand.Seed(time.Now().UnixNano())
	log.SetPrefix("[EMOGEN] ")
	log.Printf("Starting up...")
	log.Printf("Prepared %d emojis\n", len(emojis))
}

func setupRedis() redis.Pool {
	return redis.Pool{
		MaxIdle:   8,
		MaxActive: 16,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				log.Fatalf("Could not connect to redis (%s)\n", err)
			}
			return c, err
		},
	}
}

func addr() string {
	ip := os.Getenv("EMOGEN_LOCALHOST_IP")
	if ip == "" {
		ip = "127.0.0.1"
	}

	port := os.Getenv("EMOGEN_LOCALHOST_PORT")
	if port == "" {
		port = "80"
	}

	return ip + ":" + port
}

func main() {
	setup()

	redisPool := setupRedis()
	defer redisPool.Close()

	router := gin.Default()
	router.Use(cors.Default())

	setupRouter(router, &redisPool)

	router.Run(addr())
}
