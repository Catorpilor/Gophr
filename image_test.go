package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImageCreateFromURLInvalidStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	image := Image{}
	err := image.CreateFromURL(server.URL)
	if err != errInvalidImageURL {
		t.Errorf("Expected errImageURLInvalid but got %s", err)
	}
}

func TestImageCreateFromURLInvalidContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Add("Content-type", "text/html")
	}))
	defer server.Close()

	image := Image{}
	err := image.CreateFromURL(server.URL)
	if err != errInvalidImageType {
		t.Errorf("Expected errImageURLInvalid but got %s", err)
	}
}
