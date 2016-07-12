package permissions

// Permissions model for file
type Permissions struct {
	Deny         []string `json:"deny"`
	AccessUsers  []string `json:"access_users"`
	AclError     string   `json:"acl_error"`
	Acl          string   `json:"acl"`
	AccessGroups []string `json:"access_groups"`
	Allow        []string `json:"allow"`
}
