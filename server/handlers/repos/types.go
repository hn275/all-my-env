package repos

// types
type Repository struct {
	// added attributes:
	Linked          bool  `json:"linked"`
	IsOwner         bool  `json:"is_owner"`
	VariableCounter uint8 `json:"variable_counter"`

	ID          uint32  `json:"id"`
	NodeID      string  `json:"node_id"`
	Name        string  `json:"name"`
	FullName    string  `json:"full_name"`
	Private     bool    `json:"private"`
	Owner       Owner   `json:"owner"`
	HTMLURL     string  `json:"html_url"`
	Description *string `json:"description"`
	Fork        bool    `json:"fork"`
	URL         string  `json:"url"`
}

type Owner struct {
	Login      string `json:"login"`
	ID         uint32 `json:"id"`
	NodeID     string `json:"node_id"`
	AvatarURL  string `json:"avatar_url"`
	GravatarID string `json:"gravatar_id"`
	URL        string `json:"url"`
	HTMLURL    string `json:"html_url"`
}
