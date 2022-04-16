//go:build !solution

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func read(r io.Reader) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- scan.Text()
		}
	}()
	return lines
}

func main() {
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	c, _, err := websocket.DefaultDialer.Dial(*addr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				//log.Println("read:", err)
				return
			}
			fmt.Print(string(message))
		}
	}()

	//scanner := bufio.NewScanner(os.Stdin)
	mes := read(os.Stdin)
	for {
		select {
		case <-done:
			return
		case <-stop:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				//log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Millisecond * 100):
			}
			return
		case text := <-mes:
			err := c.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Println("write:", err)
				return
			}

		}
	}
}
