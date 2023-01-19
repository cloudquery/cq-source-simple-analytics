package simpleanalytics

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// AllExportFields includes all standard fields listed on https://docs.simpleanalytics.com/api/export-data-points
var AllExportFields = []string{
	"added_unix",
	"added_iso",
	"hostname",
	"hostname_original",
	"path",
	"query",
	"is_unique",
	"is_robot",
	"document_referrer",
	"utm_source",
	"utm_medium",
	"utm_campaign",
	"utm_content",
	"utm_term",
	"scrolled_percentage",
	"duration_seconds",
	"viewport_width",
	"viewport_height",
	"screen_width",
	"screen_height",
	"user_agent",
	"device_type",
	"country_code",
	"browser_name",
	"browser_version",
	"os_name",
	"os_version",
	"lang_region",
	"lang_language",
	"uuid",
}

type DataPoint struct {
	AddedUnix          uint64                 `json:"added_unix"`
	AddedISO           time.Time              `json:"added_iso"`
	Hostname           string                 `json:"hostname"`
	HostnameOriginal   string                 `json:"hostname_original"`
	Path               string                 `json:"path"`
	Query              string                 `json:"query"`
	IsUnique           bool                   `json:"is_unique"`
	IsRobot            bool                   `json:"is_robot"`
	DocumentReferrer   string                 `json:"document_referrer"`
	UTMSource          string                 `json:"utm_source"`
	UTMMedium          string                 `json:"utm_medium"`
	UTMCampaign        string                 `json:"utm_campaign"`
	UTMContent         string                 `json:"utm_content"`
	UTMTerm            string                 `json:"utm_term"`
	ScrolledPercentage float64                `json:"scrolled_percentage"`
	DurationSeconds    float64                `json:"duration_seconds"`
	ViewportWidth      int64                  `json:"viewport_width"`
	ViewportHeight     int64                  `json:"viewport_height"`
	ScreenWidth        int64                  `json:"screen_width"`
	ScreenHeight       int64                  `json:"screen_height"`
	UserAgent          string                 `json:"user_agent"`
	DeviceType         string                 `json:"device_type"`
	CountryCode        string                 `json:"country_code"`
	BrowserName        string                 `json:"browser_name"`
	BrowserVersion     string                 `json:"browser_version"`
	OSName             string                 `json:"os_name"`
	OSVersion          string                 `json:"os_version"`
	LangRegion         string                 `json:"lang_region"`
	LangLanguage       string                 `json:"lang_language"`
	UUID               string                 `json:"uuid"`
	Metadata           map[string]interface{} `json:"-"`
}

// ExportOptions sets options for the export method
type ExportOptions struct {
	Hostname string
	Start    time.Time
	End      time.Time
	Fields   []string
}

// Export returns all data points for the given time range
func (c *Client) Export(ctx context.Context, opts ExportOptions, out chan<- DataPoint) error {
	if len(opts.Fields) == 0 {
		opts.Fields = AllExportFields
	}
	values := url.Values{}
	values.Set("start", opts.Start.Format(time.RFC3339))
	values.Set("end", opts.End.Format(time.RFC3339))
	values.Set("fields", strings.Join(opts.Fields, ","))
	values.Set("version", "5")
	values.Set("format", "ndjson")
	values.Set("hostname", opts.Hostname)

	reader, err := c.get(ctx, "/api/export/datapoints", values)
	if err != nil {
		return fmt.Errorf("failed to export data points: %w", err)
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var v DataPoint
		b := scanner.Bytes()
		if err := json.Unmarshal(b, &v); err != nil {
			return fmt.Errorf("failed to decode JSON: %w", err)
		}
		m := map[string]interface{}{}
		if err := json.Unmarshal(b, &m); err != nil {
			return fmt.Errorf("failed to decode metadata fields in JSON: %w", err)
		}
		v.Metadata = map[string]interface{}{}
		for k, mv := range m {
			if strings.HasPrefix(k, "metadata.") && mv != nil {
				v.Metadata[k[9:]] = mv
			}
		}
		out <- v
	}
	return nil
}
