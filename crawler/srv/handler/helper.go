package handler

import (
	scanner "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/crawler/srv/scan/fake"
	"github.com/kazoup/platform/crawler/srv/scan/local"
	"github.com/micro/go-micro/errors"
)

func mapScanner(id int64, typ string, conf map[string]string) (scanner.Scanner, error) {
	var err error
	var s scanner.Scanner

	switch typ {
	case "fake":
		s = fake.NewFake(id, conf)
	case "local":
		s = local.NewLocal(id, "/home", conf)
	default:
		err = errors.BadRequest("go.micro.srv.crawler.Crawl.Start", "Error creating scanner")
	}

	return s, err
}
