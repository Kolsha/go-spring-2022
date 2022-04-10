//go:build !solution

package main

import (
	"flag"
	"gitlab.com/slon/shad-go/firewall/cmd/firewall/internal/firewall"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	confPath     = flag.String("conf", "./configs/example.yaml", "path to config file")
	firewallAddr = flag.String("addr", ":8081", "address of firewall")
	serviceAddr  = flag.String("service-addr", "http://eu.httpbin.org/", "address of protected service")
)

func main() {
	flag.Parse()
	confData, err := ioutil.ReadFile(*confPath)
	if err != nil {
		log.Fatal(err)
	}

	rules, err := firewall.ParseRules(confData)
	if err != nil {
		log.Fatal(err)
	}

	parsedUrl, err := url.Parse(*serviceAddr)
	if err != nil {
		log.Fatal(err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(parsedUrl)
	reverseProxy.Transport = firewall.NewFirewall(rules)

	http.HandleFunc("/", reverseProxy.ServeHTTP)
	err = http.ListenAndServe(*firewallAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
