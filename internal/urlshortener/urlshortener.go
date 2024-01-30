package urlshortener

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"strings"
)

// urlShortener represents the URL shortening application.
type urlShortener struct {
	storage map[string]string
}

// New creates a new instance of URLShortener.
func New() *urlShortener {
	return &urlShortener{storage: make(map[string]string)}
}

// Handler handles incoming HTTP requests.
func (u *urlShortener) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		u.homeHandler(w, r)
	case "/shorten":
		u.shortenHandler(w, r)
	case "/show":
		u.showHandler(w, r)
	default:
		u.redirectHandler(w, r)
	}
}

// homeHandler handles the home page.
func (u *urlShortener) homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to URL Shortener!")
}

// shortenHandler handles URL shortening.
func (u *urlShortener) shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	longURL := r.Form.Get("url")
	if longURL == "" {
		http.Error(w, "Missing 'url' parameter", http.StatusBadRequest)
		return
	}

	shortURL := u.shortenURL(longURL)
	u.storage[shortURL] = longURL

	response := map[string]string{"short_url": shortURL}
	u.writeJSONResponse(w, response)
}

func (u *urlShortener) showHandler(w http.ResponseWriter, r *http.Request) {
	response := u.storage
	u.writeJSONResponse(w, response)
}

// redirectHandler handles URL redirection.
func (u *urlShortener) redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	params := strings.Split(r.URL.Path, "/")
	shortURL := params[len(params)-1]

	longURL, ok := u.storage[shortURL]
	if !ok {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusSeeOther)
}

// shortenURL shortens the given URL.
func (u *urlShortener) shortenURL(longURL string) string {
	h := fnv.New32a()
	h.Write([]byte(longURL))
	hash := h.Sum32()

	var chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var base62Builder strings.Builder

	for hash > 0 {
		base62Builder.WriteByte(chars[hash%62])
		hash /= 62
	}

	return base62Builder.String()
}

// writeJSONResponse writes JSON response to the given HTTP response writer.
func (u *urlShortener) writeJSONResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
