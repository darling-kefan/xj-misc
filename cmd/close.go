package main

import (
	"fmt"
	"time"
)

type T int

func IsClosed(ch <-chan T) bool {
	select {
	case _, ok := <-ch:
		if !ok {
			return true
		}
	default:
	}

	return false
}

func main() {
	c := make(chan T)
	go func() {
		c <- 110
	}()
	fmt.Println(IsClosed(c)) // false
	time.Sleep(1e9)
	fmt.Println(IsClosed(c)) // false

	close(c)
	fmt.Println(IsClosed(c))
}
