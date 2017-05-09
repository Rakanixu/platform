package mock

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/lib/db/bulk"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type mock struct{}

func init() {
	bulk.Register(new(mock))
}

func (e *mock) Init(srv micro.Service) error {
	return nil
}

func (e *mock) Files(ctx context.Context, msg *crawler.FileMessage) error {
	return nil
}

func (e *mock) SlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	return nil
}

func (e *mock) SlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	return nil
}
