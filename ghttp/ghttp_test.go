package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/victorguidi/ghttp/ghttp"
)

var (
	post = map[string]string{"name": "John Wick"}
	put  = map[string]any{"name": "John Wick", "age": 30}
)

func TestHandler(t *testing.T) {
	// Setup
	e := ghttp.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.Context = *ghttp.NewContext(rec, req)
	test(e.Context)

	type Response struct {
		Name string `json:"name"`
	}
	var resp Response

	res := rec.Result()
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if resp.Name != "John Wick" {
		t.Errorf("Expected John Wick, got %s", resp.Name)
	}
}

func TestPathHandler(t *testing.T) {
	// Setup
	e := ghttp.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("name", "stich")
	rec := httptest.NewRecorder()
	e.Context = *ghttp.NewContext(rec, req)
	testParam(e.Context)

	type Response struct {
		Name string `json:"name"`
	}
	var resp Response

	res := rec.Result()
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if resp.Name != "stich" {
		t.Errorf("Expected stich, got %s", resp.Name)
	}
}

func TestPostHandler(t *testing.T) {
	payload, _ := json.Marshal(post)
	e := ghttp.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(payload))
	rec := httptest.NewRecorder()
	e.Context = *ghttp.NewContext(rec, req)
	testPost(e.Context)

	type Response struct {
		Name string `json:"name"`
	}
	var resp Response

	res := rec.Result()
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if resp.Name != "John Wick" {
		t.Errorf("Expected John Wick, got %s", resp.Name)
	}
}

func TestPutHandler(t *testing.T) {
	payload, _ := json.Marshal(put)
	e := ghttp.New()
	req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(payload))
	req.SetPathValue("name", "John")
	rec := httptest.NewRecorder()
	e.Context = *ghttp.NewContext(rec, req)
	testPut(e.Context)

	type Response struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var resp Response

	res := rec.Result()
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Errorf("expected error to be nil got %v", err)
	}
	if resp.Name != "John" || resp.Age != 30 {
		t.Errorf("Expected Name to be John and Age to be 30, got %s, %d", resp.Name, resp.Age)
	}
}

func TestDeleteHandler(t *testing.T) {
	// Setup
	e := ghttp.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("name", "stich")
	rec := httptest.NewRecorder()
	e.Context = *ghttp.NewContext(rec, req)
	testDelete(e.Context)

	if rec.Result().StatusCode != 200 {
		t.Errorf("Something crashed")
	}
}
