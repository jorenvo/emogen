package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
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

func setupRouter(r *gin.Engine) {
	r.GET("/r/:link", func(c *gin.Context) {
		link := c.Param("link")
		fmt.Printf("link is %s\n", link)

		// TODO: lookup link, if found 301, if not 404
		c.Redirect(http.StatusMovedPermanently, "https://google.com/")
	})

	r.POST("/r", func(c *gin.Context) {
		link := c.PostForm("link")
		fmt.Printf("link is %s\n", link)

		var increment uint = 1295200259 // TODO explain

		// TODO temp
		var currentEmojiNumber uint = 0
		currentEmojiNumber += increment

		emojiNumberMax := uint(math.Pow(float64(len(emojis)), 3))

		shortLink := fmt.Sprintf("/r/%s", getEmojis(currentEmojiNumber, uint(len(emojis))))

		nextEmojiNumber := getNextEmojiNumber(emojiNumberMax, increment, currentEmojiNumber)
		currentEmojiNumber = nextEmojiNumber // todo temp

		c.JSON(200, gin.H{
			"link": shortLink,
		})
	})
}

func setup() {
	// I could shuffle emojis.go manually but I'm too lazy. The
	// constant seed ensures the shuffling always produces the same
	// result.
	const mySeed = 42861292382474619
	rand.Seed(mySeed)
	rand.Shuffle(len(emojis), func(i, j int) { emojis[i], emojis[j] = emojis[j], emojis[i] })

	fmt.Printf("Prepared %d emojis\n", len(emojis))
}

func main() {
	setup()
	r := gin.Default()
	setupRouter(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
