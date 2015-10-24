package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestUsage(t *testing.T) {
	if usage("Test") != "Usage: Test [USERNAME]\n" {
		t.Error("Wrong usage message")
	}
}

func TestGetApiUrl(t *testing.T) {
	url := getApiUrl("Test")
	if url != "https://api.github.com/users/test/keys" {
		t.Errorf("Wrong api url: %s", url)
	}
}

func TestGetGithubKeysConnectionError(t *testing.T) {
	_, err := getGithubKeys("htt://wringurl")
	if err == nil {
		t.Error("Missing API error result")
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

func TestGetGithubKeysServer(t *testing.T) {
	testCases := []struct {
		s              fakeServer
		ExpectError    bool
		ExpectedResult githubKeys
	}{
		{
			s:           fakeServer{http.StatusNotFound, `{"message": "Error"}`},
			ExpectError: true,
		},
		{
			s:           fakeServer{http.StatusNotFound, `[{wrong}]`},
			ExpectError: true,
		},
		{
			s:           fakeServer{http.StatusOK, `[{wrong}]`},
			ExpectError: true,
		},
		{
			s:              fakeServer{http.StatusOK, `[{"key":"testKey1"},{"key":"testKey2"}]`},
			ExpectError:    false,
			ExpectedResult: githubKeys{{Key: "testKey1"}, {Key: "testKey2"}},
		},
	}

	for _, c := range testCases {
		testServer := httptest.NewServer(c.s)
		res, err := getGithubKeys(testServer.URL)
		if c.ExpectError && err == nil {
			t.Errorf("Missing error at case: %v", c.s)
		} else if !c.ExpectError {
			if err != nil {
				t.Errorf("Not expected error %v, at case: %v", err, c.s)
			}
			if !reflect.DeepEqual(c.ExpectedResult, res) {
				t.Errorf("Wrong result! Expected: %s, Got: %s\n", c.ExpectedResult, res)
			}
		}
	}
}
