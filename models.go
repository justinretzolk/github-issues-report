package main

type Data struct {
	Reviewed bool
	NodeId   string
	Number   int
	Title    string
	Url      string
}

type IssuesQuery struct {
	Repository struct {
		Issues struct {
			Nodes    []issue
			PageInfo pageinfo
		} `graphql:"issues(filterBy: {states: OPEN, since: $start}, first: 100, after: $cursor)"`
	} `graphql:"repository(name:$repository, owner: $owner)"`
	RateLimit ratelimit
}

type issue struct {
	Number int
	Title  string
	Url    string
}

type pageinfo struct {
	EndCursor   string
	HasNextPage bool
}

type ratelimit struct {
	Cost      int
	Remaining int
	ResetAt   string
}

type githubAuthFile struct {
	OAuthToken string `yaml:"oauth_token"`
}
