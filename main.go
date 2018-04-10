package main

import (
	"fmt"
)

func main() {

	ch := make(chan string)

	defer func() {

		fmt.Println("world")

	}()
	go func() {
		ch <- "hello"
		close(ch)
	}()
	fmt.Println(<-ch)
}
