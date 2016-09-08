package handler

import (
	"strings"

	scanner "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/crawler/srv/scan/fake"
	"github.com/kazoup/platform/crawler/srv/scan/googledrive"
	"github.com/kazoup/platform/crawler/srv/scan/local"
	"github.com/kazoup/platform/crawler/srv/scan/onedrive"
	"github.com/kazoup/platform/crawler/srv/scan/slack"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/errors"
)

func MapScanner(id int64, dataSource *datasource.Endpoint) (scanner.Scanner, error) {
	var err error
	var s scanner.Scanner
	var config map[string]string

	dsUrl := strings.Split(dataSource.Url, ":")

	switch dsUrl[0] {
	case globals.Fake:
		s = fake.NewFake(id, config)
	case globals.Local:
		s, err = local.NewLocal(id, dsUrl[1], dataSource.Index, config)
	case globals.Slack:
		s = slack.NewSlack(id, dataSource)
	case globals.GoogleDrive:
		s = googledrive.NewGoogleDrive(id, dataSource)
	case globals.OneDrive:
		s = onedrive.NewOneDrive(id, dataSource)
	default:
		err = errors.BadRequest("go.micro.srv.crawler.Crawl.Start", "Error creating scanner")
	}

	return s, err
}
