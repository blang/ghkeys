package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type githubKeys []struct {
	Key string `json:"key"`
}

type githubError struct {
	Message string `json:"message"`
}

func usage(arg string) string {
	return fmt.Sprintf("Usage: %s [USERNAME]\n", arg)
}

func getApiUrl(user string) string {
	return fmt.Sprintf("https://api.github.com/users/%s/keys", strings.ToLower(user))
}

func getGithubKeys(url string) (githubKeys, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var ghError githubError
		err = dec.Decode(&ghError)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(ghError.Message)
	}
	var ghKeys githubKeys
	err = dec.Decode(&ghKeys)
	if err != nil {
		return nil, err
	}

	return ghKeys, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage(os.Args[0]))
		os.Exit(2)
	}
	user := strings.Join(os.Args[1:], " ")
	if user == "" {
		fmt.Fprintf(os.Stderr, "Error: Invalid username\n")
		fmt.Fprintf(os.Stderr, usage(os.Args[0]))
		os.Exit(2)
	}

	ghKeys, err := getGithubKeys(getApiUrl(user))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error requesting api: %s\n", err)
		os.Exit(1)
	}

	for _, key := range ghKeys {
		fmt.Printf("%s\n", key.Key)
	}
}
