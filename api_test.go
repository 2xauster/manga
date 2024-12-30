package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type TestRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

func TestReadJSON(t *testing.T) {
	validBody := `{"name": "John Doe", "email": "john.doe@example.com"}`
	invalidBody := `{"name": "Jo", "email": "invalid"}`

	t.Run("Valid JSON", func(t *testing.T) {
		var target TestRequest
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		
		if err := readJSON(req, &target); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		var target TestRequest
		req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(invalidBody))
		req.Header.Set("Content-Type", "application/json")
		
		if err := readJSON(req, &target); err == nil {
			t.Fatal("Expected error, got nil")
		}
	})
}

func TestWriteJSON(t *testing.T) {
	response := Response{
		Status:     200,
		StatusText: "OK",
		D:          map[string]string{"message": "Success"},
		Time:       time.Now(),
	}

	w := httptest.NewRecorder()

	if err := writeJSON(response, w); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, w.Result().StatusCode)
	}

	var resBody Response
	if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resBody.Status != response.Status {
		t.Errorf("Expected status %v, got %v", response.Status, resBody.Status)
	}
	if resBody.StatusText != response.StatusText {
		t.Errorf("Expected status text %v, got %v", response.StatusText, resBody.StatusText)
	}
}

func TestHandleRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) Response {
		return newResponse(
			http.StatusOK,
			map[string]string{"message": "Test successful"},
			nil,
		)
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	h := handleRequest(handler)
	h.ServeHTTP(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %v, got %v", http.StatusOK, w.Result().StatusCode)
	}

	var resBody Response
	if err := json.NewDecoder(w.Body).Decode(&resBody); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if resBody.StatusText != "OK" {
		t.Errorf("Expected status text 'OK', got %v", resBody.StatusText)
	}
	if resBody.D.(map[string]interface{})["message"] != "Test successful" {
		t.Errorf("Expected message 'Test successful', got %v", resBody.D)
	}
}

