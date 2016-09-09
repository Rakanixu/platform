package file

type KazoupLocalFile struct {
	KazoupFile
}

func (kf *KazoupLocalFile) PreviewURL() string {
	return kf.URL
}
