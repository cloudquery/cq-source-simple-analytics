package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudquery/cq-source-simple-analytics/internal/simpleanalytics"
	"github.com/cloudquery/plugin-sdk/backend"
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/cloudquery/plugin-sdk/specs"
	"github.com/rs/zerolog"
)

type Client struct {
	Logger   zerolog.Logger
	SAClient *simpleanalytics.Client
	Backend  backend.Backend
	Website  WebsiteSpec
	websites []WebsiteSpec
}

func (c *Client) ID() string {
	return strings.Join([]string{"simple-analytics", c.Website.Hostname}, ":")
}

func (c *Client) withWebsite(website WebsiteSpec) *Client {
	return &Client{
		Logger:   c.Logger.With().Str("hostname", website.Hostname).Logger(),
		SAClient: c.SAClient,
		Backend:  c.Backend,
		Website:  website,
		websites: c.websites,
	}
}

func New(_ context.Context, logger zerolog.Logger, s specs.Source, opts source.Options) (schema.ClientMeta, error) {
	var pluginSpec Spec
	if err := s.UnmarshalSpec(&pluginSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plugin spec: %w", err)
	}
	err := pluginSpec.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate plugin spec: %w", err)
	}

	saClient := simpleanalytics.NewClient(pluginSpec.UserId, pluginSpec.APIKey)
	return &Client{
		Logger:   logger,
		Backend:  opts.Backend,
		SAClient: saClient,
		websites: pluginSpec.Websites,
	}, nil
}
