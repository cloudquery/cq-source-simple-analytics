package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudquery/cq-source-simple-analytics/client"
	"github.com/cloudquery/cq-source-simple-analytics/internal/simpleanalytics"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/transformers"
	"golang.org/x/sync/errgroup"
)

const tableDataPoints = "simple_analytics_data_points"

func DataPoints() *schema.Table {
	return &schema.Table{
		Name:        tableDataPoints,
		Description: "https://docs.simpleanalytics.com/api/export-data-points",
		Resolver:    fetchDataPoints,
		Multiplex:   client.WebsiteMultiplex,
		Transform: transformers.TransformWithStruct(
			&simpleanalytics.DataPoint{},
			transformers.WithPrimaryKeys("UUID"),
		),
		Columns: []schema.Column{
			{
				Name:     "metadata",
				Type:     schema.TypeJSON,
				Resolver: schema.PathResolver("Metadata"),
			},
		},
		IsIncremental: true,
	}
}

func fetchDataPoints(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	c := meta.(*client.Client)

	// Fetch cursor, but default to year SA was founded otherwise (assuming there is no data before that)
	start := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	if c.Backend != nil {
		value, err := c.Backend.Get(ctx, tableDataPoints, c.ID())
		if err != nil {
			return fmt.Errorf("failed to get cursor from backend: %w", err)
		}
		if value != "" {
			c.Logger.Info().Str("cursor", value).Msg("cursor found")
			start, err = time.Parse(time.RFC3339, value)
			if err != nil {
				return fmt.Errorf("failed to parse cursor from backend: %w", err)
			}
		}
	}
	end := time.Now()
	c.Logger.Info().Time("start", start).Time("end", end).Msg("fetching data points")

	// Stream data points from Simple Analytics, from start time to now.
	fields := make([]string, len(simpleanalytics.AllExportFields))
	copy(fields, simpleanalytics.AllExportFields)
	for _, field := range c.Website.MetadataFields {
		fields = append(fields, "metadata."+field)
	}
	opts := simpleanalytics.ExportOptions{
		Hostname: c.Website.Hostname,
		Start:    start,
		End:      end,
		Fields:   fields,
	}
	g, gctx := errgroup.WithContext(ctx)
	var ch = make(chan simpleanalytics.DataPoint)
	g.Go(func() error {
		defer close(ch)
		return c.SAClient.Export(gctx, opts, ch)
	})
	for range ch {
		res <- <-ch
	}
	err := g.Wait()
	if err != nil {
		return fmt.Errorf("failed to fetch data points: %w", err)
	}

	// Save cursor state to the backend.
	if c.Backend != nil {
		// We subtract 15 minutes from the end time to allow for delayed data points
		// to be fetched on the next sync. This will cause some duplicates, but
		// allows us to guarantee at-least-once delivery. Duplicates can be removed
		// by using overwrite-delete-stale write mode, by de-duplicating in queries,
		// or by running a post-processing step.
		newCursor := end.Add(-15 * time.Minute).Format(time.RFC3339)
		err = c.Backend.Set(ctx, tableDataPoints, c.ID(), newCursor)
		if err != nil {
			return fmt.Errorf("failed to save cursor to backend: %w", err)
		}
		c.Logger.Info().Str("cursor", newCursor).Msg("cursor updated")
	}
	return nil
}
