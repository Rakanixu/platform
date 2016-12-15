package gmail

type GmailFile struct {
	Id           string `json:"id"`
	MessageId    string `json:"message_id"`
	Name         string `json:"name"`
	MimeType     string `json:"mime_type"`
	Extension    string `json:"extension"`
	Base64       string `json:"base64"`
	InternalDate int64  `json:"internalDate,omitempty,string"`
	SizeEstimate int64  `json:"sizeEstimate,omitempty"`
}
