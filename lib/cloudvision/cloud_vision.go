package cloudvision

import (
	"cloud.google.com/go/vision"
	"golang.org/x/net/context"
	"io"
)

func Tag(rd io.ReadCloser) ([]string, error) {
	var s []string

	defer rd.Close()

	client, err := vision.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	image, err := vision.NewImageFromReader(rd)
	if err != nil {
		return nil, err
	}

	labels, err := client.DetectLabels(context.Background(), image, 10)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		s = append(s, label.Description)
	}

	return s, nil
}
