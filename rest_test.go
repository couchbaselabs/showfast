package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code: %v", resp.StatusCode)
	}
}

func TestAdmin(t *testing.T) {
	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/admin")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code: %v", resp.StatusCode)
	}
}

func TestStatic(t *testing.T) {
	ts := httptest.NewServer(httpEngine())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/static/js/timeline.js")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("unexpected status code: %v", resp.StatusCode)
	}
}
