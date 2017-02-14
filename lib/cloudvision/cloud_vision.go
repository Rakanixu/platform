package cloudvision

import (
	"cloud.google.com/go/vision"
	"fmt"
	"golang.org/x/net/context"
	"io"
)

func Tag(rd io.ReadCloser) error {
	defer rd.Close()

	client, err := vision.NewClient(context.Background())
	if err != nil {
		return err
	}

	image, err := vision.NewImageFromReader(rd)
	if err != nil {
		return err
	}

	labels, err := client.DetectLabels(context.Background(), image, 10)
	if err != nil {
		return err
	}

	fmt.Println("Labels:")
	for _, label := range labels {
		fmt.Println(label.Description)
	}

	return nil
}
