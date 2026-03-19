package main

import (
	"demo/or"
	"fmt"
	"time"
)

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	var ch1 chan interface{}
	var ch12 chan interface{}
	select {
	case <-or.Or(
		ch1, ch12,
	):
	case <-time.After(100 * time.Millisecond):
		fmt.Println("all")
		return
	}

	fmt.Printf("done after %v\n", time.Since(start))
}
