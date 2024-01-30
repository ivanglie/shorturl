package main

import (
	"log"
	"net/http"

	"github.com/ivanglie/shorturl/internal/urlshortener"
)

func main() {
	urlShortener := urlshortener.New()

	log.Println("Server is running on :8080")
	http.HandleFunc("/", urlShortener.Handler)
	http.ListenAndServe(":8080", nil)
}
