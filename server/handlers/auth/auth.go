package auth

import (
	"net/http"
)

type authCtx interface {
	Do(r *http.Request) (*http.Response, error)
}

var (
	authClient authCtx
)

func init() {
	authClient = &http.Client{}
}

type GithubUser struct {
	// for db
	Login string `json:"login"`
	ID    uint32 `json:"id"`
	Email string `json:"email"`
	// for display
	Bio               string `json:"bio"`
	Name              string `json:"name"`
	AvatarURL         string `json:"avatar_url"`
	HTMLURL           string `json:"html_url"`
	PublicRepos       uint16 `json:"public_repos"`
	TotalPrivateRepos uint16 `json:"total_private_repos"`
	OwnedPrivateRepos uint16 `json:"owned_private_repos"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	/* unused fields
	NodeID                  string      `json:"node_id"`
	GravatarID              string      `json:"gravatar_id"`
	URL                     string      `json:"url"`
	FollowersURL            string      `json:"followers_url"`
	FollowingURL            string      `json:"following_url"`
	GistsURL                string      `json:"gists_url"`
	StarredURL              string      `json:"starred_url"`
	SubscriptionsURL        string      `json:"subscriptions_url"`
	OrganizationsURL        string      `json:"organizations_url"`
	ReposURL                string      `json:"repos_url"`
	EventsURL               string      `json:"events_url"`
	ReceivedEventsURL       string      `json:"received_events_url"`
	Type                    string      `json:"type"`
	SiteAdmin               bool        `json:"site_admin"`
	Company                 interface{} `json:"company"`
	Blog                    string      `json:"blog"`
	Location                interface{} `json:"location"`
	Hireable                interface{} `json:"hireable"`
	TwitterUsername         interface{} `json:"twitter_username"`
	PublicGists             int64       `json:"public_gists"`
	Followers               int64       `json:"followers"`
	Following               int64       `json:"following"`
	PrivateGists            int64       `json:"private_gists"`
	DiskUsage               int64       `json:"disk_usage"`
	Collaborators           int64       `json:"collaborators"`
	TwoFactorAuthentication bool        `json:"two_factor_authentication"`
	Plan                    Plan        `json:"plan"`
	*/
}

type Plan struct {
	Name          string `json:"name"`
	Space         int64  `json:"space"`
	Collaborators int64  `json:"collaborators"`
	PrivateRepos  int64  `json:"private_repos"`
}
