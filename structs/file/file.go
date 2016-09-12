package file

const (
	DEFAULT_IMAGE_PREVIEW_URL string = "http://www.scaleautomag.com/sitefiles/images/no-preview-available.png"
)

type File interface {
	PreviewURL() string
}
