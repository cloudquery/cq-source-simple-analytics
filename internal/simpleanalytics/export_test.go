package simpleanalytics

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/sync/errgroup"
)

const (
	testHostname = "saasforcovid.com"
	testUserID   = "sa_user_id_77969473-8121-4ef4-882b-2bda8acc7fc3"
	testAPIKey   = "sa_api_key_xwPSzcqDIjb4xNZVM76WYMb3LNCbstdkmttT"
)

func TestExportPageViews(t *testing.T) {
	// read contents of testdata file
	f := "testdata/pageviews.ndjson"
	contents, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("unexpected error reading testdata file: %v", err)
	}

	var gotRequest *http.Request
	// create a test server that returns the contents of the testdata file
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequest = r
		w.WriteHeader(http.StatusOK)
		w.Write(contents)
	}))

	defer ts.Close()
	c := NewClient(testUserID, testAPIKey, WithBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	start := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	end := time.Now().AddDate(0, 0, 0).Truncate(24 * time.Hour)
	opts := ExportOptions{
		Hostname: testHostname,
		Start:    start,
		End:      end,
		Fields:   []string{},
	}
	got := testExportPageViews(t, c, opts)

	if gotRequest == nil {
		t.Fatalf("expected request to test server, got none")
	}
	if gotRequest.URL.Path != "/api/export/datapoints" {
		t.Errorf("unexpected path in request. got: %s, want: %s", gotRequest.URL.Path, "/api/export/datapoints")
	}
	q := gotRequest.URL.Query()
	if q.Get("hostname") != testHostname {
		t.Errorf("unexpected hostname in request. got: %s, want: %s", q.Get("hostname"), testHostname)
	}
	if q.Get("format") != "ndjson" {
		t.Errorf("unexpected format in request. got: %s, want: %s", q.Get("format"), "ndjson")
	}
	if q.Get("start") != opts.Start.Format("2006-01-02") {
		t.Errorf("unexpected start in request. got: %s, want: %s", q.Get("start"), opts.Start.Format("2006-01-02"))
	}
	if q.Get("end") != opts.End.Format("2006-01-02") {
		t.Errorf("unexpected end in request. got: %s, want: %s", q.Get("end"), opts.End.Format("2006-01-02"))
	}

	// last UUID in testdata file
	if got[len(got)-1].UUID != "0bec40c0-06a6-43b2-ac38-b209a08de836" {
		t.Errorf("unexpected last UUID in results. got: %s, want: %s", got[len(got)-1].UUID, "0bec40c0-06a6-43b2-ac38-b209a08de836")
	}
	wantMetadata := map[string]any{
		"fieldname_text": "test",
		"fieldname_date": "2023-01-20T14:57:59.698Z",
		"fieldname_bool": true,
		"fieldname_int":  123.0,
	}
	if diff := cmp.Diff(wantMetadata, got[len(got)-1].Metadata); diff != "" {
		t.Errorf("unexpected metadata in last result. diff: %s", diff)
	}
}

func TestExportPageViewsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	c := NewClient(testUserID, testAPIKey)
	start := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	end := time.Now()
	opts := ExportOptions{
		Hostname: testHostname,
		Start:    start,
		End:      end,
	}
	testExportPageViews(t, c, opts)
}

func testExportPageViews(t *testing.T, c *Client, opts ExportOptions) []PageView {
	t.Helper()
	results := make(chan PageView)
	g := errgroup.Group{}
	g.Go(func() error {
		defer close(results)
		return c.ExportPageViews(context.Background(), opts, results)
	})
	got := make([]PageView, 0)
	for r := range results {
		got = append(got, r)
	}
	err := g.Wait()
	if err != nil {
		t.Fatalf("unexpected error calling ExportPageViews: %v", err)
	}

	if len(got) == 0 {
		t.Fatalf("expected at least one result, got 0")
	}
	for _, r := range got {
		if r.UUID == "" {
			t.Fatalf("unexpected empty UUID. Full row: %v", r)
		}
	}
	return got
}

func TestExportEvents(t *testing.T) {
	// read contents of testdata file
	f := "testdata/events.ndjson"
	contents, err := os.ReadFile(f)
	if err != nil {
		t.Fatalf("unexpected error reading testdata file: %v", err)
	}

	var gotRequest *http.Request
	// create a test server that returns the contents of the testdata file
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequest = r
		w.WriteHeader(http.StatusOK)
		w.Write(contents)
	}))

	defer ts.Close()
	c := NewClient(testUserID, testAPIKey, WithBaseURL(ts.URL), WithHTTPClient(ts.Client()))
	start := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	end := time.Now().AddDate(0, 0, 0).Truncate(24 * time.Hour)
	opts := ExportOptions{
		Hostname: testHostname,
		Start:    start,
		End:      end,
		Fields:   []string{},
	}
	got := testExportEvents(t, c, opts)

	if gotRequest == nil {
		t.Fatalf("expected request to test server, got none")
	}
	if gotRequest.URL.Path != "/api/export/datapoints" {
		t.Errorf("unexpected path in request. got: %s, want: %s", gotRequest.URL.Path, "/api/export/datapoints")
	}
	q := gotRequest.URL.Query()
	if q.Get("hostname") != testHostname {
		t.Errorf("unexpected hostname in request. got: %s, want: %s", q.Get("hostname"), testHostname)
	}
	if q.Get("format") != "ndjson" {
		t.Errorf("unexpected format in request. got: %s, want: %s", q.Get("format"), "ndjson")
	}
	if q.Get("start") != opts.Start.Format("2006-01-02") {
		t.Errorf("unexpected start in request. got: %s, want: %s", q.Get("start"), opts.Start.Format("2006-01-02"))
	}
	if q.Get("end") != opts.End.Format("2006-01-02") {
		t.Errorf("unexpected end in request. got: %s, want: %s", q.Get("end"), opts.End.Format("2006-01-02"))
	}

	// last added_iso in testdata file
	pt, _ := time.Parse(time.RFC3339, "2023-01-23T12:43:35.689Z")
	if !got[len(got)-1].AddedISO.Equal(pt) {
		t.Errorf("unexpected added_iso in results. got: %s, want: %s", got[len(got)-1].AddedISO, pt)
	}
	wantMetadata := map[string]any{
		"fieldname_text": "test",
		"fieldname_date": "2023-01-20T14:57:59.698Z",
		"fieldname_bool": true,
		"fieldname_int":  123.0,
	}
	if diff := cmp.Diff(wantMetadata, got[len(got)-1].Metadata); diff != "" {
		t.Errorf("unexpected metadata in last result. diff: %s", diff)
	}
}

func TestExportEventsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	c := NewClient(testUserID, testAPIKey)
	start := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	end := time.Now()
	opts := ExportOptions{
		Hostname: testHostname,
		Start:    start,
		End:      end,
		Fields: append(ExportFieldsEvents, []string{
			"metadata.updated_date",
			"metadata.boolean_bool",
			"metadata.decimalCount_int",
			"metadata.string_text",
		}...),
	}
	testExportEvents(t, c, opts)
}

func testExportEvents(t *testing.T, c *Client, opts ExportOptions) []Event {
	t.Helper()
	results := make(chan Event)
	g := errgroup.Group{}
	g.Go(func() error {
		defer close(results)
		return c.ExportEvents(context.Background(), opts, results)
	})
	got := make([]Event, 0)
	for r := range results {
		got = append(got, r)
	}
	err := g.Wait()
	if err != nil {
		t.Fatalf("unexpected error calling ExportEvents: %v", err)
	}

	if len(got) == 0 {
		t.Fatalf("expected at least one result, got 0")
	}
	for _, r := range got {
		if r.Hostname == "" {
			t.Fatalf("unexpected empty hostname. Full row: %v", r)
		}
	}
	return got
}
