package dropbox

import (
	"time"
)

type FilesListResponse struct {
	Entries []DropboxFile `json:"entries"`
	Cursor  string        `json:"cursor"`
	HasMore bool          `json:"has_more"`
}

type DropboxFile struct {
	//Tag            string    `json:".tag"`
	Name           string    `json:"name"`
	PathLower      string    `json:"path_lower"`
	PathDisplay    string    `json:"path_display"`
	ID             string    `json:"id"`
	ClientModified time.Time `json:"client_modified"`
	ServerModified time.Time `json:"server_modified"`
	Rev            string    `json:"rev"`
	Size           int       `json:"size"`
	MediaInfo      struct {
		//Tag      string `json:".tag"`
		Metadata struct {
			//Tag        string `json:".tag"`
			Dimensions struct {
				Height int `json:"height"`
				Width  int `json:"width"`
			} `json:"dimensions"`
		} `json:"metadata"`
	} `json:"media_info"`
	HasExplicitSharedMembers bool `json:"has_explicit_shared_members"`
}

/*
// I remove just the Tag property, so I can keep going quickly and not have to implement the following TODO:


type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// Fucking retarded developers..
//https://www.dropboxforum.com/hc/en-us/community/posts/204403136-API-v2-feedback-files-tag-attribute-extremely-unfriendly-for-HTTP-JSON
func (df DropboxFile) MarshalJSON() ([]byte, error) {
	//TODO: do a proper marshaller, so ".tag" is just "tag" , therefore elastic search won't trow the exception
	tag, _ := json.Marshal(df.Tag)
	name, _ := json.Marshal(df.Name)
	//path_lower, _ := json.Marshal(df.PathLower)
	//path_display, _ := json.Marshal(df.PathDisplay)

	return []byte(`{
		"tag":` + string(tag) + `,
		"name":` + string(name) + `}`)
}
*/
