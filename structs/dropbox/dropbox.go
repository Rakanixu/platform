package dropbox

import (
	"time"
)

type FilesListResponse struct {
	Entries []DropboxFile `json:"entries"`
	Cursor  string        `json:"cursor"`
	HasMore bool          `json:"has_more"`
}

// DropboxFile represent a dropbox file
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
	HasExplicitSharedMembers bool             `json:"has_explicit_shared_members"`
	DropboxUsers             []DropboxUser    `json:"dropbox_users"`   // This field is calculated in different request, only if file is shared
	DropboxInvitees          []DropboxInvitee `json:"dropbox_invitee"` // This field is calculated in different request, only if file is shared
}

type FileMembersListResponse struct {
	Users []struct { // We just want the accountId to be able to retrieve a DropboxUser
		User struct {
			AccountID string `json:"account_id"`
		} `json:"user"`
	} `json:"users"`
	Groups   []DropboxGroup   `json:"groups"`
	Invitees []DropboxInvitee `json:"invitees"`
}

// DropboxUser represent a dropbox user / account
type DropboxUser struct {
	AccountID string `json:"account_id"`
	Name      struct {
		GivenName       string `json:"given_name"`
		Surname         string `json:"surname"`
		FamiliarName    string `json:"familiar_name"`
		DisplayName     string `json:"display_name"`
		AbbreviatedName string `json:"abbreviated_name"`
	} `json:"name"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	Disabled        bool   `json:"disabled"`
	IsTeammate      bool   `json:"is_teammate"`
	ProfilePhotoURL string `json:"profile_photo_url"`
}

// DropboxGroup represent a dropbox group. Do not know how it works, so it is not implemented
type DropboxGroup struct {
	AccessType struct {
		//Tag string `json:".tag"`
	} `json:"access_type"`
	Group struct {
		GroupName           string `json:"group_name"`
		GroupID             string `json:"group_id"`
		GroupManagementType struct {
			//Tag string `json:".tag"`
		} `json:"group_management_type"`
		GroupType struct {
			//Tag string `json:".tag"`
		} `json:"group_type"`
		IsOwner     bool `json:"is_owner"`
		SameTeam    bool `json:"same_team"`
		MemberCount int  `json:"member_count"`
	} `json:"group"`
	Permissions []interface{} `json:"permissions"`
	IsInherited bool          `json:"is_inherited"`
}

// DropboxInvitee represent a user that has access to a dropbox file, but does not have a dropbox account
type DropboxInvitee struct {
	Permissions []interface{} `json:"permissions"`
	IsInherited bool          `json:"is_inherited"`
	Invitee     struct {
		//Tag   string `json:".tag"`
		Email string `json:"email"`
	} `json:"invitee"`
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
