package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloword(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(helloWorld)

	handler.ServeHTTP(rr, req)

	expected_status := http.StatusOK
	if rr.Code != expected_status {
		t.Fatalf("handler returned wrong status code: got %v want %v", rr.Code, expected_status)
	}

	expected_text := "Hello World"
	if rr.Body.String() != expected_text {
		t.Fatalf("handler returned wrong status code: got %v want %v",rr.Code, expected_text)
	}
}