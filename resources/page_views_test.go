package resources

import (
	"testing"

	"github.com/cloudquery/cq-source-simple-analytics/client"
)

func TestPageViews(t *testing.T) {
	ts := testServer(t, "page_views.ndjson")
	defer ts.Close()
	client.TestHelper(t, PageViews(), ts)
}
