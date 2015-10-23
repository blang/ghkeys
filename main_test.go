package main

import (
	"bytes"
	"io"
	"os"
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
