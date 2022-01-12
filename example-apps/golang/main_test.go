package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, rr.Code, expected_status, "handler returned wrong status code: got %v want %v", rr.Code, expected_status)
	
	expected_text := "Hello World"
	assert.HTTPBodyContains(t, helloWorld, "GET", "/", nil, expected_text ,"handler returned wrong status code: got %v want %v",rr.Code, expected_text)
}