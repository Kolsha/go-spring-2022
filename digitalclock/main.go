//go:build !solution

package main

import (
	"flag"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func main() {
	port := flag.String("port", "80", "http server port")
	flag.Parse()

	http.HandleFunc("/", getPngHandler)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}

func getPngHandler(w http.ResponseWriter, r *http.Request) {
	var sizeParam = 1
	var timeParam string
	var err error

	queries := r.URL.Query()

	if sizeSlc, ok := queries["k"]; ok && len(sizeSlc) > 0 {
		sizeParam, err = strconv.Atoi(sizeSlc[0])
		if err != nil || sizeParam < 1 || sizeParam > 30 {
			http.Error(w, "invalid k param", http.StatusBadRequest)
			return
		}
	}

	if timeSlc, ok := queries["time"]; ok {
		timeParam = timeSlc[0]
		timePattern := "^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$"
		ok, err = regexp.MatchString(timePattern, timeParam)

		if err != nil || !ok {
			http.Error(w, "invalid time param", http.StatusBadRequest)
			return
		}
	} else {
		t := time.Now()
		timeParam = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
	}

	fmt.Println(sizeParam, timeParam)
	img := CreateImage(timeParam, sizeParam)

	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, img)

	if err != nil {
		http.Error(w, "render image error", http.StatusBadRequest)
	}
}
