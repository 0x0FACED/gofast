package main

import (
	"os/exec"
	"testing"
	"time"
)

func BenchmarkFindLinux(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("find", "/", "-name", "*image*.png")
		start := time.Now()

		output, err := cmd.CombinedOutput()
		if err != nil {
			b.Logf("Error executing find: %v\nOutput: %s", err, string(output))
		}

		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Milliseconds()), "ms/op")
	}
}

func BenchmarkGofast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("go", "run", "main.go", "-p", "/", "-t", "file", "-n", "*image*.png", "-m", "pattern")
		start := time.Now()

		out, err := cmd.CombinedOutput()
		if err != nil {
			b.Fatalf("Error executing gofast: %v\nOut: %s", err, string(out))
		}

		b.Logf("Output: %s", string(out))

		elapsed := time.Since(start)
		b.ReportMetric(float64(elapsed.Milliseconds()), "ms/op")
	}
}
