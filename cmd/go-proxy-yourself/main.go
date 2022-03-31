package main

import (
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		host, _, _ := net.SplitHostPort(req.Host)

		if host == "mfa.nneisen.local" {
			io.WriteString(w, host)
		}
	}

	http.HandleFunc("/", helloHandler)
	log.Println("Listing for requests")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
