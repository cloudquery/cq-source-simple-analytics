package client

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// DefaultStartTime defaults to the year SA was founded (we assume there were no data before that)
var DefaultStartTime = time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)

// AllowedTimeLayout is the layout used for the start_time and end_time fields, and matches what the export API supports
var AllowedTimeLayout = "2006-01-02"

var reValidDuration = regexp.MustCompile(`^(\d+)([dmy])$`)

type Spec struct {
	// UserID is the Simple Analytics API user ID.
	UserID string `json:"user_id"`

	// APIKey is the Simple Analytics API key.
	APIKey string `json:"api_key"`

	// Websites is a list of websites to fetch data for.
	Websites []WebsiteSpec `json:"websites"`

	// StartDateStr is the time to start fetching data from. If specified, it must use AllowedTimeLayout.
	StartDateStr string `json:"start_date"`

	// EndDateStr is the time at which to stop fetching data. If not specified, the current time is used.
	// If specified, it must use AllowedTimeLayout.
	EndDateStr string `json:"end_date"`

	// PeriodStr is the duration of the time window to fetch historical data for, in days, months or years.
	// Examples:
	//  "7d": past 7 days
	//  "3m": last 3 months
	//  "1y": last year
	// It is used to calculate start_time if it is not specified. If start_time is specified,
	// duration is ignored.
	PeriodStr string `json:"duration"`
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
	if s.StartDateStr != "" {
		_, err := time.Parse(AllowedTimeLayout, s.StartDateStr)
		if err != nil {
			return fmt.Errorf("could not parse start_time: %v", err)
		}
	}
	if s.EndDateStr != "" {
		_, err := time.Parse(AllowedTimeLayout, s.EndDateStr)
		if err != nil {
			return fmt.Errorf("could not parse end_time: %v", err)
		}
	}
	if s.PeriodStr != "" {
		_, err := parsePeriod(s.PeriodStr)
		if err != nil {
			return fmt.Errorf("could not validate period: %v (should be a number followed by \"d\", \"m\" or \"y\", e.g. \"7d\", \"1m\" or \"3y\")", err)
		}
	}
	return nil
}

func (s *Spec) SetDefaults() {
	if s.StartDateStr == "" && s.PeriodStr == "" {
		s.StartDateStr = DefaultStartTime.Format(AllowedTimeLayout)
	}
	if s.EndDateStr == "" {
		s.EndDateStr = time.Now().Format(AllowedTimeLayout)
	}
}

func (s Spec) StartTime() time.Time {
	if s.StartDateStr == "" && s.PeriodStr != "" {
		return time.Now().Add(-s.Period())
	}
	t, _ := time.Parse(AllowedTimeLayout, s.StartDateStr) // any error should be caught by Validate()
	return t
}

func (s Spec) EndTime() time.Time {
	t, _ := time.Parse(AllowedTimeLayout, s.EndDateStr) // any error should be caught by Validate()
	return t
}

func (s Spec) Period() time.Duration {
	d, _ := parsePeriod(s.PeriodStr) // any error should be caught by Validate()
	return d
}

func parsePeriod(s string) (time.Duration, error) {
	m := reValidDuration.FindStringSubmatch(s)
	if m == nil {
		return 0, errors.New("invalid period")
	}
	n, err := strconv.Atoi(m[1])
	if err != nil {
		return 0, err
	}
	switch m[2] {
	case "d":
		return time.Duration(n) * 24 * time.Hour, nil
	case "m":
		return time.Duration(n) * 30 * 24 * time.Hour, nil
	case "y":
		return time.Duration(n) * 365 * 24 * time.Hour, nil
	}
	// should never happen, we already validated using regex
	panic("unhandled period unit")
}
