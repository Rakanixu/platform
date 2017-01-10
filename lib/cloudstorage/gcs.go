package cloudstorage

import "io"

type GoogleCloudStorage struct {
}

func (gcs *GoogleCloudStorage) Upload(r io.Reader, s string) error {
	return nil
}
