package file

type KazoupLocalFile struct {
	KazoupFile
}

func (kf *KazoupLocalFile) PreviewURL() string {
	return DEFAULT_IMAGE_PREVIEW_URL
}
