package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudquery/cq-source-simple-analytics/client"
	"github.com/cloudquery/cq-source-simple-analytics/internal/simpleanalytics"
	"github.com/cloudquery/plugin-sdk/faker"
)

func TestEvents(t *testing.T) {
	var pv simpleanalytics.Event
	if err := faker.FakeObject(&pv); err != nil {
		t.Fatal(err)
	}
	pv.Metadata = map[string]any{
		"metadata.foo_text": "bar",
		"metadata.bar_int":  123,
		"metadata.baz_date": "2021-01-01T00:00:00Z",
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ndjson")
		d, _ := json.Marshal(pv)
		_, _ = w.Write(d)
	}))
	defer ts.Close()
	client.TestHelper(t, Events(), ts)
}
