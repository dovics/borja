package main

import (
	"log"
	"net/http"
)

func main() {
	resp, err := http.Get("http://127.0.0.1:8081/clear")
	if err != nil {
		log.Fatal("http request error: ", err)
	}

	if resp.StatusCode != 200 {
		log.Fatal("http request error")
	}
}
