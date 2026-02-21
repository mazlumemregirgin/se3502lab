package lab

import (
	"fmt"
	"math/rand"
	"time"
)

func Guess() {
	fmt.Println("Welcome to the Guessing Game!")
	fmt.Println("I have selected a random number between 1 and 100. Can you guess it?")

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	target := random.Intn(100) + 1

	var tahmin int
	for {
		fmt.Print("Enter your tahmin: ")

		_, err := fmt.Scanln(&tahmin)
		if err != nil {
			fmt.Println("Lütfen geçerli bir sayı giriniz.")
			continue
		}

		if tahmin < target {
			fmt.Println("Too low!")
		} else if tahmin > target {
			fmt.Println("Too high!")
		} else {
			fmt.Println("You got it!")
			break
		}
	}
}
