package main

import (
	"fmt"
	"math/rand"
	"time"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("think of a number between 1 and 100")
	num := rand.Intn(100)
	count := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter guess: ")
		text, _ := reader.ReadString('\n')
		if guess, err := strconv.Atoi(strings.TrimSpace(text)); err != nil {
			fmt.Println("Not a Number!!")
		} else {
			count++
			if guess > num {
				fmt.Println("Too High")
			} else if guess < num {
				fmt.Println("Too Low")
			} else {
				break;
			}
		}
	}
	fmt.Printf("You got it in %d", count)

}