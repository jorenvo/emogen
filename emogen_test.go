package main

// import "fmt"
import (
	"math"
	"testing"
)

func testIncrementOverNumbers(t *testing.T, increment uint, currentNumber uint, numbers []bool) {
	for i := 0; i < len(numbers); i++ {
		currentNumber = getNextEmojiNumber(uint(len(numbers)), increment, currentNumber)

		if numbers[currentNumber] {
			t.Errorf("%d was already in the sequence", currentNumber)
			break
		}

		numbers[currentNumber] = true
	}
}

func TestGetNextEmojiEven(t *testing.T) {
	numbers := make([]bool, 128)
	var currentNumber uint = 0
	var increment uint = 41

	testIncrementOverNumbers(t, increment, currentNumber, numbers)
}

func TestGetNextEmojiOdd(t *testing.T) {
	numbers := make([]bool, 129)
	var currentNumber uint = 0
	var increment uint = 71

	testIncrementOverNumbers(t, increment, currentNumber, numbers)
}

func TestGetEmojis1(t *testing.T) {
	// For a 9 bit number we would have 3 groups of 3 bits each. The
	// max number in each group is 2^3.
	var length uint = uint(math.Pow(float64(2), float64(3)))
	var number uint = 5 << 6 | 6 << 3 | 7

	n1, n2, n3 := getEmojiNumbers(number, length)

	if n1 != 7 {
		t.Errorf("n1 is %d, expected 7", n1)
	}

	if n2 != 6 {
		t.Errorf("n2 is %d, expected 6", n2)
	}

	if n3 != 5 {
		t.Errorf("n3 is %d, expected 5", n3)
	}
}

func TestGetEmojis2(t *testing.T) {
	// For a 9 bit number we would have 3 groups of 3 bits each. The
	// max number in each group is 2^3.
	var length uint = uint(math.Pow(float64(2), float64(3)))
	var number uint = 7 << 6 | 7 << 3 | 7

	n1, n2, n3 := getEmojiNumbers(number, length)

	if n1 != 7 {
		t.Errorf("n1 is %d, expected 7", n1)
	}

	if n2 != 7 {
		t.Errorf("n2 is %d, expected 7", n2)
	}

	if n3 != 7 {
		t.Errorf("n3 is %d, expected 7", n3)
	}
}
