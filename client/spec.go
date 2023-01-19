package client

import (
	"fmt"
	"time"
)

// DefaultStartTime defaults to the year SA was founded (we assume there were no data before that)
var DefaultStartTime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

// AllowedTimeLayout is the layout used for the start_time and end_time fields, and matches what the export API supports
var AllowedTimeLayout = "2006-01-02"

type Spec struct {
	// UserID is the Simple Analytics API user ID.
	UserID string `json:"user_id"`

	// APIKey is the Simple Analytics API key.
	APIKey string `json:"api_key"`

	// Websites is a list of websites to fetch data for.
	Websites []WebsiteSpec `json:"websites"`

	// StartTimeStr is the time to start fetching data from. If specified, it must use AllowedTimeLayout.
	StartTimeStr string `json:"start_time"`

	// EndTimeStr is the time at which to stop fetching data. If not specified, the current time is used.
	// If specified, it must use AllowedTimeLayout.
	EndTimeStr string `json:"end_time"`

	// WindowOverlapSeconds gives a number of seconds to decrease the start_time by
	// when starting from an incremental cursor position. This allows for late-arriving data to
	// be fetched in a subsequent sync and guarantee at-least-once delivery, but can
	// introduce duplicates.
	WindowOverlapSeconds int `json:"window_overlap_seconds"`
}

type WebsiteSpec struct {
	Hostname       string   `json:"hostname"`
	MetadataFields []string `json:"metadata_fields"`
}

func (s Spec) Validate() error {
	if s.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if s.APIKey == "" {
		return fmt.Errorf("api_key is required")
	}
	if len(s.Websites) == 0 {
		return fmt.Errorf("at least one website is required")
	}
	for _, w := range s.Websites {
		if w.Hostname == "" {
			return fmt.Errorf("every website entry must have a hostname")
		}
	}
	if s.StartTimeStr != "" {
		_, err := time.Parse(AllowedTimeLayout, s.StartTimeStr)
		if err != nil {
			return fmt.Errorf("could not parse start_time: %v", err)
		}
	}
	if s.EndTimeStr != "" {
		_, err := time.Parse(AllowedTimeLayout, s.EndTimeStr)
		if err != nil {
			return fmt.Errorf("could not parse end_time: %v", err)
		}
	}
	return nil
}

func (s *Spec) SetDefaults() {
	if s.StartTimeStr == "" {
		s.StartTimeStr = DefaultStartTime.Format(AllowedTimeLayout)
	}
	if s.EndTimeStr == "" {
		s.EndTimeStr = time.Now().Format(AllowedTimeLayout)
	}
	if s.WindowOverlapSeconds == 0 {
		s.WindowOverlapSeconds = 60
	}
}

func (s Spec) StartTime() time.Time {
	t, _ := time.Parse(AllowedTimeLayout, s.StartTimeStr) // any error should be caught by Validate()
	return t
}

func (s Spec) EndTime() time.Time {
	t, _ := time.Parse(AllowedTimeLayout, s.EndTimeStr) // any error should be caught by Validate()
	return t
}
