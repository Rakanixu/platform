package dropbox

import (
	"encoding/json"
	"time"
)

type FilesListResponse struct {
	Entries []DropboxFile `json:"entries"`
	Cursor  string        `json:"cursor"`
	HasMore bool          `json:"has_more"`
}

// DropboxFile represent a dropbox file
// Unmarshall using default UnmarshalJSON(We get all info internally, including Tag string `json:".tag"` )
// Custom Marshall to not include Tag string `json:".tag"`, so elastic search does not trow mapping exceptions
type DropboxFile struct {
	Tag                      string           `json:".tag"`
	Name                     string           `json:"name"`
	PathLower                string           `json:"path_lower"`
	PathDisplay              string           `json:"path_display"`
	ID                       string           `json:"id"`
	ClientModified           time.Time        `json:"client_modified"`
	ServerModified           time.Time        `json:"server_modified"`
	Rev                      string           `json:"rev"`
	Size                     int              `json:"size"`
	MediaInfo                MediaInfo        `json:"media_info"`
	HasExplicitSharedMembers bool             `json:"has_explicit_shared_members"`
	PublicURL                string           `json:"public_url"`      // This field is calculated in different request, only if file is public to evryone
	DropboxTag               string           `json:"dropbox_tag"`     // This field does not exists on response, but we will fill constructing KazoupDropboxFile
	DropboxUsers             []DropboxUser    `json:"dropbox_users"`   // This field is calculated in different request, only if file is shared
	DropboxInvitees          []DropboxInvitee `json:"dropbox_invitee"` // This field is calculated in different request, only if file is shared
}

type PublicFilesListResponse struct {
	Links   []DropboxPublicFile `json:"links"`
	Cursor  string              `json:"cursor"`
	HasMore bool                `json:"has_more"`
}

type DropboxPublicFile struct {
	URL       string `json:"url"`
	Name      string `json:"name"`
	PathLower string `json:"path_lower"`
	ID        string `json:"id"`
}

type MediaInfo struct {
	Metadata Metadata `json:"metadata"`
}

type Metadata struct {
	Dimensions Dimensions `json:"dimensions"`
}

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
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

// Fucking retarded developers.. elastic search will throw mappijng exceptions cause they do not follow JSON standard
// https://www.dropboxforum.com/hc/en-us/community/posts/204403136-API-v2-feedback-files-tag-attribute-extremely-unfriendly-for-HTTP-JSON
func (df DropboxFile) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		// Tag, this is the field we do not want to marshal to not be push to ES
		Name                     string           `json:"name"`
		PathLower                string           `json:"path_lower"`
		PathDisplay              string           `json:"path_display"`
		ID                       string           `json:"id"`
		ClientModified           time.Time        `json:"client_modified"`
		ServerModified           time.Time        `json:"server_modified"`
		Rev                      string           `json:"rev"`
		Size                     int              `json:"size"`
		MediaInfo                MediaInfo        `json:"media_info"`
		HasExplicitSharedMembers bool             `json:"has_explicit_shared_members"`
		PublicURL                string           `json:"public_url"`      // This field is calculated in different request, only if file is public to evryone
		DropboxTag               string           `json:"dropbox_tag"`     // This field does not exists on response, but we will fill constructing KazoupDropboxFile
		DropboxUsers             []DropboxUser    `json:"dropbox_users"`   // This field is calculated in different request, only if file is shared
		DropboxInvitees          []DropboxInvitee `json:"dropbox_invitee"` // This field is calculated in different request, only if file is shared
	}{
		Name:                     df.Name,
		PathLower:                df.PathLower,
		PathDisplay:              df.PathDisplay,
		ID:                       df.ID,
		ClientModified:           df.ClientModified,
		ServerModified:           df.ServerModified,
		Rev:                      df.Rev,
		Size:                     df.Size,
		MediaInfo:                df.MediaInfo,
		HasExplicitSharedMembers: df.HasExplicitSharedMembers,
		PublicURL:                df.PublicURL,
		DropboxTag:               df.DropboxTag,
		DropboxUsers:             df.DropboxUsers,
		DropboxInvitees:          df.DropboxInvitees,
	})
}
