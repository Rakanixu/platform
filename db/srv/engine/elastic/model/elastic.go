package model

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
)

type FilesChannel struct {
	FileMessage *crawler.FileMessage
	Client      client.Client
	Ctx         context.Context
}

type Elastic struct {
	Client               *elib.Client
	BulkProcessor        *elib.BulkProcessor
	BulkFilesProcessor   *elib.BulkProcessor
	FilesChannel         chan *FilesChannel
	SlackUsersChannel    chan *crawler.SlackUserMessage
	SlackChannelsChannel chan *crawler.SlackChannelMessage
}
