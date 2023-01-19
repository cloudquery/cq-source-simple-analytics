package client

import "github.com/cloudquery/plugin-sdk/schema"

func WebsiteMultiplex(meta schema.ClientMeta) []schema.ClientMeta {
	var l = make([]schema.ClientMeta, 0)
	client := meta.(*Client)
	for _, website := range client.websites {
		l = append(l, client.withWebsite(website))
	}
	return l
}
