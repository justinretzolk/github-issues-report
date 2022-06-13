package main

type Data struct {
	Reviewed bool
	NodeId   string
	Number   int
	Title    string
	Url      string
}

type IssuesQuery struct   {
	Search struct {
		PageInfo pageinfo
		Nodes []issue_fragment
	} `graphql:"search(type:$searchType, first:100, after:$cursor, query:$searchQuery)"`
	RateLimit ratelimit
}

type issue_fragment struct {
	Issue issue `graphql:"...on Issue"`
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
