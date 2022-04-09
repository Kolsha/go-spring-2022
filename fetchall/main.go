//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func fetch(url string, ch chan<- time.Duration) {
	start := time.Now()

	resp, err := http.Get(url)

	if err == nil {
		defer resp.Body.Close()
	}
	ch <- time.Since(start)
}

func main() {
	start := time.Now()

	ch := make(chan time.Duration)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	close(ch)

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
