package client

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cloudquery/cq-source-simple-analytics/internal/simpleanalytics"
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

func TestHelper(t *testing.T, table *schema.Table, ts *httptest.Server) {
	version := "vDev"
	table.IgnoreInTests = false
	t.Helper()
	l := zerolog.New(zerolog.NewTestWriter(t)).Output(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.StampMicro},
	).Level(zerolog.DebugLevel).With().Timestamp().Logger()
	newTestExecutionClient := func(ctx context.Context, logger zerolog.Logger, spec specs.Source, opts source.Options) (schema.ClientMeta, error) {
		saClient := simpleanalytics.NewClient("test", "test", simpleanalytics.WithBaseURL(ts.URL), simpleanalytics.WithHTTPClient(ts.Client()))
		s := Spec{
			UserID: "test",
			APIKey: "test",
			Websites: []WebsiteSpec{
				{
					Hostname:       "test.com",
					MetadataFields: []string{"metadata_text", "metadata_int"},
				},
			},
		}
		s.SetDefaults()
		err := s.Validate()
		if err != nil {
			return nil, err
		}
		return &Client{
			Logger:   l,
			SAClient: saClient,
			Backend:  opts.Backend,
			Spec:     s,
		}, nil
	}
	p := source.NewPlugin(
		table.Name,
		version,
		[]*schema.Table{
			table,
		},
		newTestExecutionClient)
	p.SetLogger(l)
	source.TestPluginSync(t, p, specs.Source{
		Name:         "dev",
		Path:         "cloudquery/dev",
		Version:      version,
		Tables:       []string{table.Name},
		Destinations: []string{"mock-destination"},
	})
}
