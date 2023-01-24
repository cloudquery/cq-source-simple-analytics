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

const dateLayout = "2006-01-02"

var ExportFieldsEvents = []string{
	"added_iso",
	"added_unix",
	"browser_name",
	"browser_version",
	"country_code",
	"datapoint",
	"device_type",
	"document_referrer",
	"hostname",
	"hostname_original",
	"is_robot",
	"lang_language",
	"lang_region",
	"os_name",
	"os_version",
	"path",
	"path_and_query",
	"query",
	"screen_height",
	"screen_width",
	"session_id",
	"user_agent",
	"utm_campaign",
	"utm_content",
	"utm_medium",
	"utm_source",
	"utm_term",
	"uuid",
	"viewport_height",
	"viewport_width",
}

var ExportFieldsPageViews = []string{
	"added_iso",
	"added_unix",
	"browser_name",
	"browser_version",
	"country_code",
	"device_type",
	"document_referrer",
	"duration_seconds",
	"hostname",
	"hostname_original",
	"is_robot",
	"is_unique",
	"lang_language",
	"lang_region",
	"os_name",
	"os_version",
	"path",
	"path_and_query",
	"query",
	"screen_height",
	"screen_width",
	"scrolled_percentage",
	"session_id",
	"user_agent",
	"utm_campaign",
	"utm_content",
	"utm_medium",
	"utm_source",
	"utm_term",
	"uuid",
	"viewport_height",
	"viewport_width",
}

type Event struct {
	AddedISO         time.Time              `json:"added_iso"`
	AddedUnix        uint64                 `json:"added_unix"`
	BrowserName      string                 `json:"browser_name"`
	BrowserVersion   string                 `json:"browser_version"`
	CountryCode      string                 `json:"country_code"`
	Datapoint        string                 `json:"datapoint"`
	DeviceType       string                 `json:"device_type"`
	DocumentReferrer string                 `json:"document_referrer"`
	Hostname         string                 `json:"hostname"`
	HostnameOriginal string                 `json:"hostname_original"`
	IsRobot          bool                   `json:"is_robot"`
	LangLanguage     string                 `json:"lang_language"`
	LangRegion       string                 `json:"lang_region"`
	Metadata         map[string]interface{} `json:"-"`
	OSName           string                 `json:"os_name"`
	OSVersion        string                 `json:"os_version"`
	Path             string                 `json:"path"`
	PathAndQuery     string                 `json:"path_and_query"`
	Query            string                 `json:"query"`
	ScreenHeight     int64                  `json:"screen_height"`
	ScreenWidth      int64                  `json:"screen_width"`
	SessionID        string                 `json:"session_id"`
	UTMCampaign      string                 `json:"utm_campaign"`
	UTMContent       string                 `json:"utm_content"`
	UTMMedium        string                 `json:"utm_medium"`
	UTMSource        string                 `json:"utm_source"`
	UTMTerm          string                 `json:"utm_term"`
	UUID             string                 `json:"uuid"`
	UserAgent        string                 `json:"user_agent"`
	ViewportHeight   int64                  `json:"viewport_height"`
	ViewportWidth    int64                  `json:"viewport_width"`
}

type PageView struct {
	AddedISO           time.Time              `json:"added_iso"`
	AddedUnix          uint64                 `json:"added_unix"`
	BrowserName        string                 `json:"browser_name"`
	BrowserVersion     string                 `json:"browser_version"`
	CountryCode        string                 `json:"country_code"`
	DeviceType         string                 `json:"device_type"`
	DocumentReferrer   string                 `json:"document_referrer"`
	DurationSeconds    float64                `json:"duration_seconds"`
	Hostname           string                 `json:"hostname"`
	HostnameOriginal   string                 `json:"hostname_original"`
	IsRobot            bool                   `json:"is_robot"`
	IsUnique           bool                   `json:"is_unique"`
	LangLanguage       string                 `json:"lang_language"`
	LangRegion         string                 `json:"lang_region"`
	Metadata           map[string]interface{} `json:"-"`
	OSName             string                 `json:"os_name"`
	OSVersion          string                 `json:"os_version"`
	Path               string                 `json:"path"`
	PathAndQuery       string                 `json:"path_and_query"`
	Query              string                 `json:"query"`
	ScreenHeight       int64                  `json:"screen_height"`
	ScreenWidth        int64                  `json:"screen_width"`
	ScrolledPercentage float64                `json:"scrolled_percentage"`
	SessionID          string                 `json:"session_id"`
	UTMCampaign        string                 `json:"utm_campaign"`
	UTMContent         string                 `json:"utm_content"`
	UTMMedium          string                 `json:"utm_medium"`
	UTMSource          string                 `json:"utm_source"`
	UTMTerm            string                 `json:"utm_term"`
	UUID               string                 `json:"uuid"`
	UserAgent          string                 `json:"user_agent"`
	ViewportHeight     int64                  `json:"viewport_height"`
	ViewportWidth      int64                  `json:"viewport_width"`
}

// ExportOptions sets options for the export method
type ExportOptions struct {
	Hostname string
	Start    time.Time
	End      time.Time
	Fields   []string
}

// ExportPageViews returns all page views for the given time range
func (c *Client) ExportPageViews(ctx context.Context, opts ExportOptions, out chan<- PageView) error {
	if len(opts.Fields) == 0 {
		opts.Fields = ExportFieldsPageViews
	}
	values := getQueryParams(opts)
	values.Set("type", "pageviews")

	reader, err := c.get(ctx, "/api/export/datapoints", values)
	if err != nil {
		return fmt.Errorf("failed to export data points: %w", err)
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var v PageView
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

// ExportEvents returns all events for the given time range
func (c *Client) ExportEvents(ctx context.Context, opts ExportOptions, out chan<- Event) error {
	if len(opts.Fields) == 0 {
		opts.Fields = ExportFieldsEvents
	}
	values := getQueryParams(opts)
	values.Set("type", "events")

	reader, err := c.get(ctx, "/api/export/datapoints", values)
	if err != nil {
		return fmt.Errorf("failed to export data points: %w", err)
	}
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		var v Event
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

func getQueryParams(opts ExportOptions) url.Values {
	values := url.Values{}
	values.Set("start", opts.Start.Format(dateLayout))
	values.Set("end", opts.End.Format(dateLayout))
	values.Set("fields", strings.Join(opts.Fields, ","))
	values.Set("version", "5")
	values.Set("format", "ndjson")
	values.Set("hostname", opts.Hostname)
	return values
}
