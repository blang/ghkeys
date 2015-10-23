package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type GithubKeys []struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [USERNAME]\n")
}

func main() {

	// No username
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
		fmt.Fprintf(os.Stderr, "Error: %s\n")
		os.Exit(1)
	}
	dec := json.NewDecoder(resp.Body)
	var ghKeys GithubKeys
	err = dec.Decode(&ghKeys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n")
		os.Exit(1)
	}
	for _, key := range ghKeys {
		fmt.Printf("%s\n", key.Key)
	}
}
