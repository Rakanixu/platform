package handler

import (
	scanner "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/crawler/srv/scan/fake"
	"github.com/kazoup/platform/crawler/srv/scan/local"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/errors"
	"strings"
)

func mapScanner(id int64, dataSource *datasource.Endpoint) (scanner.Scanner, error) {
	var err error
	var s scanner.Scanner
	var config map[string]string

	dsUrl := strings.Split(dataSource.Url, ":")

	switch dsUrl[0] {
	case "fake":
		s = fake.NewFake(id, config)
	case "local":
		s, err = local.NewLocal(id, dsUrl[1], config)
	default:
		err = errors.BadRequest("go.micro.srv.crawler.Crawl.Start", "Error creating scanner")
	}

	return s, err
}
