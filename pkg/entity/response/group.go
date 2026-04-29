package response

// RrIdPGroup is the API representation of a group managed by the identity provider.
type RrIdPGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
}

// RrIdPGroups wraps a list of IdP groups.
type RrIdPGroups struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Groups  []RrIdPGroup `json:"groups,omitempty"`
}

// RrSingleIdPGroup wraps a single IdP group.
type RrSingleIdPGroup struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Group   *RrIdPGroup `json:"group,omitempty"`
}
