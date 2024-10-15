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

	// Test the function with the test server URL
	statusCode, err := callURL(server.URL)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, statusCode)
	}
}

func TestCallURL_Error(t *testing.T) {
	// Test the function with an invalid URL
	statusCode, err := callURL("http://invalid-url")

	if err == nil {
		t.Fatalf("Expected an error, got none")
	}

	if statusCode != 0 {
		t.Fatalf("Expected status code 0, got %d", statusCode)
	}
}
