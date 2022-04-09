//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counter := make(map[string]uint)

	for _, filePath := range os.Args[1:] {
		file, err := os.Open(filePath)
		if err != nil {
			continue
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			counter[line]++
		}

	}

	for line, count := range counter {
		if count < 2 {
			continue
		}
		fmt.Printf("%v\t%v\n", count, line)
	}
}
