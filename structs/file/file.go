package file

const (
	DEFAULT_IMAGE_PREVIEW_URL string = "http://www.scaleautomag.com/sitefiles/images/no-preview-available.png"
	BASE_URL_FILE_PREVIEW     string = "http://localhost:8082/media"
)

type File interface {
	PreviewURL() string
}
