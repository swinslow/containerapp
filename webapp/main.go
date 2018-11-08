package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	var WEBPORT string
	if WEBPORT = os.Getenv("WEBPORT"); WEBPORT == "" {
		WEBPORT = "3001"
	}

	http.HandleFunc("/", rootHandler)
	fmt.Println("Listening on :" + WEBPORT)
	log.Fatal(http.ListenAndServe(":"+WEBPORT, nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, path is %s\n", r.URL.Path)
}
