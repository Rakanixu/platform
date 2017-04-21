package bulk

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type DBBulker interface {
	Init(micro.Service) error
	Bulker
}

type Bulker interface {
	Files(ctx context.Context, msg *crawler.FileMessage) error
	SlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error
	SlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error
}

var (
	bulker DBBulker
)

func Register(storage DBBulker) {
	bulker = storage
}

func Init(srv micro.Service) error {
	return bulker.Init(srv)
}

func Files(ctx context.Context, msg *crawler.FileMessage) error {
	return bulker.Files(ctx, msg)
}

func SlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	return bulker.SlackUsers(ctx, msg)
}

func SlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	return bulker.SlackChannels(ctx, msg)
}
