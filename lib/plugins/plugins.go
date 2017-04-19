package plugins

import (
	// Plugins import
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/registry/nats"
	_ "github.com/micro/go-plugins/transport/tcp"

	// DB Implementation
	_ "github.com/kazoup/platform/lib/db/bulk/elastic"
	_ "github.com/kazoup/platform/lib/db/config/elastic"
	_ "github.com/kazoup/platform/lib/db/operations/elastic"
)
