package main

type GHRepo struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	HtmlURL  string `json:"html_url"`
	PushedAt string `json:"pushed_at"`
}

type GHError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

type GHContributions struct {
	Data struct {
		Viewer struct {
			ContributionsCollection struct {
				PullRequestContributionsByRepository []GHPullRequestContributionsByRepository `json:"pullRequestContributionsByRepository"`
			} `json:"contributionsCollection"`
		} `json:"viewer"`
	} `json:"data"`
}

type GHPullRequestContributionsByRepository struct {
	Repository    GHRepository   `json:"repository"`
	Contributions GHContribution `json:"contributions"`
}

type GHRepository struct {
	ID            int     `json:"id"`
	IsFork        bool    `json:"isFork"`
	NameWithOwner string  `json:"nameWithOwner"`
	Owner         GHOwner `json:"owner"`
	URL           string  `json:"url"`
}

type GHOwner struct {
	Login string `json:"login"`
}

type GHContribution struct {
	Nodes []GHNode `json:"nodes"`
}

type GHNode struct {
	PullRequest GHPullRequest `json:"pullRequest"`
}

type GHPullRequest struct {
	CreatedAt string `json:"createdAt"`
}
