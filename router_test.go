package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetServiceWithInvalidSwitchId(t *testing.T) {
	req, err := http.NewRequest("GET", "/switchstate?id=p629", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSwitchState)

	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusNotFound, rr.Code)
}

func TestGetServiceWithValidSwitchId(t *testing.T) {
	req, err := http.NewRequest("GET", "/switchstate?id=p62921", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSwitchState)

	handler.ServeHTTP(rr, req)
	checkResponseCode(t, http.StatusOK, rr.Code)
}

func TestUpdateServiceWithValidSwitchId(t *testing.T) {
	req, err := http.NewRequest("PUT", "/switchstate?id=p49873&state=disabled", nil)

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SetSwitchState)

	handler.ServeHTTP(rr, req)

	checkResponseCode(t, http.StatusOK, rr.Code)
}
