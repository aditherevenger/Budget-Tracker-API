package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMainOutput(t *testing.T) {
	// Save original stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call main
	main()

	// Capture output
	w.Close()
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}
	os.Stdout = old

	expected := "Temporary main: Budget Tracker API is under development.\n"
	if buf.String() != expected {
		t.Errorf("unexpected output: got %q, want %q", buf.String(), expected)
	}
}
