package model

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	lib "github.com/mattbaird/elastigo/lib"
	"github.com/micro/go-micro/client"
)

type FilesChannel struct {
	FileMessage *crawler.FileMessage
	Client      client.Client
}

type Elastic struct {
	Conn                 *lib.Conn
	Bulk                 *lib.BulkIndexer
	FilesChannel         chan *FilesChannel
	SlackUsersChannel    chan *crawler.SlackUserMessage
	SlackChannelsChannel chan *crawler.SlackChannelMessage
	CrawlerFinished      chan *crawler.CrawlerFinishedMessage
}
