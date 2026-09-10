package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

func TestNextConfirmTime(t *testing.T) {
	now := time.Date(2026, 9, 10, 10, 0, 0, 0, time.Local)

	tests := []struct {
		name     string
		schedule []string
		expected time.Duration
	}{
		{
			name:     "later today",
			schedule: []string{"11:00"},
			expected: time.Hour,
		},
		{
			name:     "tomorrow",
			schedule: []string{"09:00"},
			expected: 23 * time.Hour,
		},
		{
			name:     "earliest of multiple times",
			schedule: []string{"15:00", "11:00", "18:00"},
			expected: time.Hour,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := nextConfirmTime(now, test.schedule)

			if got != test.expected {
				t.Fatalf("Expected %v, got %v", test.expected, got)
			}
		})
	}
}
