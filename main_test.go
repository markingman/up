package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCallURL(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ch := make(chan Result, 1)

	// Test the function with the test server URL
	callURL(server.URL, ch)
	result := <-ch

	if result.httpCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, result.httpCode)
	}

	if !result.complete {
		t.Fatal("Expected result to be complete")
	}
}

func TestCallURL_Error(t *testing.T) {
	// Create and immediately close a server so the request reliably fails
	server := httptest.NewServer(nil)
	url := server.URL
	server.Close()

	ch := make(chan Result, 1)

	callURL(url, ch)
	result := <-ch

	if result.httpCode != 0 {
		t.Fatalf("Expected status code 0, got %d", result.httpCode)
	}

	if !result.complete {
		t.Fatal("Expected result to be complete")
	}
}
