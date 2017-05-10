package subscriber

import (
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	crawler "github.com/kazoup/platform/lib/protomsg/crawler"
	"golang.org/x/net/context"
	"testing"
)

var (
	discoveryFinished = new(DiscoveryFinished)
)

func TestPostDiscovery(t *testing.T) {
	var postDiscoveryTestData = []struct {
		ctx context.Context
		msg *crawler.CrawlerFinishedMessage
		out error
	}{
		{
			ctx,
			&crawler.CrawlerFinishedMessage{
				DatasourceId: "test",
			},
			nil,
		},
	}

	for _, tt := range postDiscoveryTestData {
		result := discoveryFinished.PostDiscovery(tt.ctx, tt.msg)
		if result != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, result)
		}
	}
}
