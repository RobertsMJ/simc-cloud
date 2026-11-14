package test

import (
	"encoding/json"
	"flag"
	"os"
	"testing"
)

var update = flag.Bool("update", false, "update expected test results")

func GoldenValueJSON[T any](t *testing.T, path string, actual T) T {
	if *update {
		js, err := json.MarshalIndent(actual, "", "  ")
		if err != nil {
			t.Fatal("failed to marshal golden value", "err", err)
		}
		err = os.WriteFile(path, js, 0644)
		if err != nil {
			t.Fatal("failed to write golden value", "err", err)
		}

		return actual
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal("failed to read golden value", "err", err)
	}
	var expected T
	err = json.Unmarshal(content, &expected)
	if err != nil {
		t.Fatal("failed to unmarshal existing golden value", "path", path, "err", err)
	}
	return expected
}

func GoldenValueRaw(t *testing.T, path string, actual []byte) []byte {

	if *update {
		err := os.WriteFile(path, actual, 0644)
		if err != nil {
			t.Fatal("failed to write golden value", "err", err)
		}
		return actual
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal("failed to read golden value", "err", err)
	}
	return content
}
