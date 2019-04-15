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

type ToShorten struct {
	Link string `json:"link" binding:"required"`
}

func setupRouter(router *gin.Engine, redisConn redis.Conn) {
	// This is a large prime. Using a prime as the increment for
	// getNextEmojiNumber will make it loop through all numbers.
	const emojiNumberIncrement = 1295200259

	emojiNumberMax := uint(math.Pow(float64(len(emojis)), 3))

	router.GET("/:link", func(c *gin.Context) {
		link := c.Param("link")
		log.Printf("resolving %s\n", link)

		link, err := redis.String(redisConn.Do("GET", "shortlink:"+link))
		if err != nil {
			link = "/notfound"
		}

		c.Redirect(http.StatusMovedPermanently, link)
	})

	router.POST("/", func(c *gin.Context) {
		var toShorten ToShorten
		if err := c.ShouldBindJSON(&toShorten); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		link := toShorten.Link
		log.Printf("shortening %s\n", link)

		shortLink, err := redis.String(redisConn.Do("GET", "link:"+link))
		if err == nil {
			c.JSON(200, gin.H{
				"link": "/" + shortLink,
			})
			return
		}

		currentEmojiNumber := getEmogenNr(redisConn)
		currentEmojiNumber = getNextEmojiNumber(emojiNumberMax, emojiNumberIncrement, currentEmojiNumber)

		shortLink = getEmojis(currentEmojiNumber, uint(len(emojis)))

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

		_, err = redisConn.Do("SET", "emogen:nr", currentEmojiNumber)
		if err != nil {
			log.Printf("Error while storing emogen:nr: %s\n", err)
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

	router := gin.Default()
	router.Use(cors.Default())

	setupRouter(router, redisConn)

	router.Run() // listen and serve on 0.0.0.0:8080
}
