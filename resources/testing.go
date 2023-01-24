package resources

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func testServer(t *testing.T, filename string) *httptest.Server {
	contents, err := os.ReadFile(path.Join("testdata", filename))
	if err != nil {
		t.Fatalf("unexpected error reading testdata file: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Write(contents)
	}))
	return ts
}
