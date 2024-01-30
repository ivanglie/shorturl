package urlshortener

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLShortener_Handler(t *testing.T) {
	t.Run("Home Handler", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		u.Handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Shorten Handler", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("url=https://example.com"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		u.Handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Show Handler", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/show", nil)
		u.Handler(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Redirect Handler", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/unknown", nil)
		u.Handler(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestURLShortener_HomeHandler(t *testing.T) {
	u := New()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	u.homeHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Welcome to URL Shortener!")
}

func TestURLShortener_ShortenHandler(t *testing.T) {
	t.Run("Valid POST Request", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("url=https://example.com"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		u.shortenHandler(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		assert.NoError(t, json.NewDecoder(w.Body).Decode(&response))
		assert.Contains(t, response, "short_url")

		assert.Equal(t, "https://example.com", u.storage[response["short_url"]])
	})

	t.Run("Invalid HTTP Method", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/shorten", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		u.shortenHandler(w, r)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Failed to Parse Form Data", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("url;https://example.com"))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		u.shortenHandler(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "Failed to parse form data")
	})

	t.Run("Missing 'url' Parameter", func(t *testing.T) {
		u := New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/shorten", nil)
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		u.shortenHandler(w, r)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestURLShortener_ShowHandler(t *testing.T) {
	u := New()
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/show", nil)

	u.showHandler(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.Nil(t, err)
	assert.Empty(t, response)
}

func TestURLShortener_RedirectHandler(t *testing.T) {
	u := New()

	u.storage["abc"] = "https://example.com"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/abc", nil)

	u.redirectHandler(w, r)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Location"))
}
