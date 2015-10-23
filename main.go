package main

import (
	"encoding/json"
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [USERNAME]\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	user := strings.Join(os.Args[1:], " ")
	if user == "" {
		fmt.Fprintf(os.Stderr, "Error: Invalid username\n")
		usage()
		os.Exit(2)
	}
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/keys", strings.ToLower(user)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error requesting api: %s\n", err)
		os.Exit(1)
	}
	dec := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var ghError githubError
		err = dec.Decode(&ghError)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown error occurred.\n")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "API Error: %s.\n", ghError.Message)
		os.Exit(1)
	}
	var ghKeys githubKeys
	err = dec.Decode(&ghKeys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	for _, key := range ghKeys {
		fmt.Printf("%s\n", key.Key)
	}
}
