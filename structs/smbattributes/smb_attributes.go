package smbattributes

type SmbAttributes struct {
	ExtendedAttributes int     `json:"extended_attributes"`
	Created            float64 `json:"created"`
	LastAttrChange     float64 `json:"last_attr_change"`
	LastWrite          float64 `json:"last_write"`
	LastAccess         float64 `json:"last_access"`
	Offline            bool    `json:"offline"`
}
