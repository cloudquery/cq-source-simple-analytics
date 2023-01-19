package plugin

import (
	"github.com/cloudquery/cq-source-simple-analytics/client"
	"github.com/cloudquery/cq-source-simple-analytics/resources"
	"github.com/cloudquery/plugin-sdk/plugins/source"
	"github.com/cloudquery/plugin-sdk/schema"
)

var (
	Version = "development"
)

func Plugin() *source.Plugin {
	return source.NewPlugin(
		"simple-analytics",
		Version,
		schema.Tables{
			resources.DataPoints(),
		},
		client.New,
	)
}
