package oauth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const GithubServer = "github"

var (
	githubConfig = oauth2.Config{
		ClientID:     "b9f6bf2a0ece96c9f97b",
		ClientSecret: "14f80ca11bef172d2f5fa884fbb2446188511607",
		Scopes:       []string{"user:email", "repo", "repo:status"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://89.221.208.88:8100/auth/github/",
	}
)
