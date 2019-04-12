package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"math"
	"math/rand"
	"net/http"
)

func getNextEmojiNumber(length, increment, current uint) uint {
	return (current + increment) % length
}

func getEmojiNumbers(number, length uint) (uint, uint, uint) {
	n1 := number % length
	n2 := (number / length) % length
	n3 := number / length / length

	return n1, n2, n3
}

func getEmojis(number, length uint) string {
	emoji1, emoji2, emoji3 := getEmojiNumbers(number, length)
	return fmt.Sprintf("%s%s%s", emojis[emoji1], emojis[emoji2], emojis[emoji3])
}

func getEmogenNr(c redis.Conn) uint {
	nr, err := redis.Int(c.Do("GET", "emogen:nr"))
	if err != nil {
		log.Printf("emogen:nr is not initialized (%s), initializing to 0", err)
		return 0
	}

	log.Printf("emogen:nr is %d", nr)
	return uint(nr)
}

func setupRouter(engine *gin.Engine, redisConn redis.Conn) {
	engine.GET("/r/:link", func(c *gin.Context) {
		link := c.Param("link")
		log.Printf("link is %s\n", link)

		// TODO: lookup link, if found 301, if not 404
		c.Redirect(http.StatusMovedPermanently, "https://google.com/")
	})

	engine.POST("/r", func(c *gin.Context) {
		link := c.PostForm("link")
		log.Printf("link is %s\n", link)

		// TODO temp
		currentEmojiNumber := getEmogenNr(redisConn)

		emojiNumberMax := uint(math.Pow(float64(len(emojis)), 3))  // TODO is constant
		const increment = 1295200259 // TODO explain
		currentEmojiNumber = getNextEmojiNumber(emojiNumberMax, increment, currentEmojiNumber)

		shortLink := fmt.Sprintf("/r/%s", getEmojis(currentEmojiNumber, uint(len(emojis))))

		_, err := redisConn.Do("SET", "emogen:nr", currentEmojiNumber)
		if err != nil {
			log.Printf("Error while storing emogen:nr: %s\n", err)
			c.JSON(500, "Failed connecting to db.")
			return
		}

		c.JSON(200, gin.H{
			"link": shortLink,
		})
	})
}

func setup() {
	log.SetPrefix("[EMOGEN] ")
	log.Printf("Starting up...")

	// I could shuffle emojis.go manually but I'm too lazy. The
	// constant seed ensures the shuffling always produces the same
	// result.
	const mySeed = 42861292382474619
	rand.Seed(mySeed)
	rand.Shuffle(len(emojis), func(i, j int) { emojis[i], emojis[j] = emojis[j], emojis[i] })

	log.Printf("Prepared %d emojis\n", len(emojis))
}

func setupRedis() redis.Conn {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Could not connect to redis (%s)\n", err)
	}

	return c
}

func main() {
	setup()

	redisConn := setupRedis()
	defer redisConn.Close()

	engine := gin.Default()

	setupRouter(engine, redisConn)

	engine.Run() // listen and serve on 0.0.0.0:8080
}
