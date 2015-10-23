package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestUsage(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stderr = w
	usage("Test")

	outChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outChan <- buf.String()
	}()
	w.Close()

	if <-outChan != "Usage: Test [USERNAME]\n" {
		t.Error("Wrong usage message")
	}
}

func TestGetApiUrl(t *testing.T) {
	url := getApiUrl("Test")
	if url != "https://api.github.com/users/test/keys" {
		t.Errorf("Wrong api url: %s", url)
	}
}

type fakeServer struct {
	StatusCode int
	Result     string
}

func (s fakeServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(s.StatusCode)
	w.Write([]byte(s.Result))
}

func TestGetGithubKeysConnectionError(t *testing.T) {
	_, err := getGithubKeys("htt://wringurl")
	if err == nil {
		t.Error("Missing API error result")
	}
}

func TestGetGithubKeysServerError(t *testing.T) {
	fs := fakeServer{404, `{"message": "Error"}`}
	testServer := httptest.NewServer(fs)

	_, err := getGithubKeys(testServer.URL)
	if err == nil {
		t.Error("Missing API error result")
	}
	if err.Error() != "Error" {
		t.Error("Wrong API error result", err)
	}
}

func TestGetGithubKeysServerSuccess(t *testing.T) {
	fs := fakeServer{http.StatusOK, `[{"key":"testKey1"},{"key":"testKey2"}]`}
	testServer := httptest.NewServer(fs)

	res, err := getGithubKeys(testServer.URL)
	if err != nil {
		t.Error("Wrong API error result", err)
	}

	expected := githubKeys{
		{Key: "testKey1"}, {Key: "testKey2"},
	}
	if !reflect.DeepEqual(expected, res) {
		t.Errorf("Wrong result! Expected: %s, Got: %s\n", expected, res)
	}
}
