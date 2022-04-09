//go:build !solution

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	store map[string]string
	mu    sync.Mutex
)

type requestBody struct {
	URL string
}

type responseBody struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func main() {
	port := flag.String("port", "80", "http server port")
	flag.Parse()

	store = make(map[string]string)

	http.HandleFunc("/shorten", createShortLinkHandler)
	http.HandleFunc("/go/", shortLinkRedirectHandler)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}

func createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	var rb requestBody
	decoder := json.NewDecoder(r.Body)
	if decoder.Decode(&rb) != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	short, err := randSeq(rb.URL)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}

	mu.Lock()
	store[short] = rb.URL
	mu.Unlock()

	data := responseBody{URL: rb.URL, Key: short}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(data)

	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func shortLinkRedirectHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	key := p[len(p)-1]

	mu.Lock()
	url, exist := store[key]
	mu.Unlock()

	if !exist {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusFound)
}

func randSeq(str string) (res string, e error) {
	algorithm := sha1.New()
	_, err := algorithm.Write([]byte(str))

	if err != nil {
		e = err
		return
	}

	res = hex.EncodeToString(algorithm.Sum(nil))
	return
}
