//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var athletes []Athlete

func main() {
	port := flag.String("port", "80", "http server port")
	jsonPath := flag.String("data", "./olympics/testdata/olympicWinners.json", "json path")
	flag.Parse()

	jsonFile, err := os.Open(*jsonPath)
	if err != nil {
		log.Fatal("read file error")
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	err = json.Unmarshal(byteValue, &athletes)

	if err != nil {
		log.Fatal("json parsing error")
	}

	http.HandleFunc("/athlete-info", AthleteHandler)
	http.HandleFunc("/top-athletes-in-sport", TopAthletesHandler)
	http.HandleFunc("/top-countries-in-year", TopCountriesHandler)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
