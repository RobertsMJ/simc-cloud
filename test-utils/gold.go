package test

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"testing"
)

var update = flag.Bool("update", false, "update expected test results")

func GoldenValue[T any](t *testing.T, path string, actual T) T {
	f, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		t.Fatal("failed to open golden file", "path", path)
	}
	defer f.Close()

	if *update {
		js, err := json.MarshalIndent(actual, "", "  ")
		if err != nil {
			t.Fatal("failed to marshal golden value")
		}
		_, err = f.Write(js)
		if err != nil {
			t.Fatal("failed to write golden value")
		}
		return actual
	}

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatal("failed to real golden value")
	}
	var expected T
	err = json.Unmarshal(content, &expected)
	if err != nil {
		t.Fatal("failed to marshal existing golden value")
	}
	return expected
}
