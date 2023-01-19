package client

import (
	"fmt"
)

type Spec struct {
	UserId   string        `json:"user_id"`
	APIKey   string        `json:"api_key"`
	Websites []WebsiteSpec `json:"websites"`
}

type WebsiteSpec struct {
	Hostname       string   `json:"hostname"`
	MetadataFields []string `json:"metadata_fields"`
}

func (s Spec) Validate() error {
	if s.UserId == "" {
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
	return nil
}
