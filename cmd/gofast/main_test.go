package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestSearchWithSingleWorker(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go",
		"-p", "/",
		"-t", "file",
		"-n", "postgresql.conf",
		"-m", "like",
		"-workers", "1",
	)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Error running command: %s\nStderr: %s", err, stderr.String())
	}

	if !strings.Contains(out.String(), "postgresql.conf") {
		t.Errorf("Expected output to contain 'postgresql.conf', got: %s", out.String())
	}
}

func TestSearchWithFiftyWorkers(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go",
		"-p", "/",
		"-t", "file",
		"-n", "postgresql.conf",
		"-m", "like",
		"-workers", "50",
	)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Error running command: %s\nStderr: %s", err, stderr.String())
	}

	if !strings.Contains(out.String(), "postgresql.conf") {
		t.Errorf("Expected output to contain 'postgresql.conf', got: %s", out.String())
	}
}
