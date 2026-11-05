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
		t.Fatal("failed to open golden file", "path", path, "err", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal("failed to close golden file", "path", path, "err", err)
		}
	}()

	if *update {
		js, err := json.MarshalIndent(actual, "", "  ")
		if err != nil {
			t.Fatal("failed to marshal golden value", "err", err)
		}
		_, err = f.Write(js)
		if err != nil {
			t.Fatal("failed to write golden value", "err", err)
		}
		return actual
	}

	content, err := io.ReadAll(f)
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
