package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"

	"github.com/HiDeoo/alfred-workflow-tools/pkg/request"
)

var (
	client     request.Client
	queryLimit = 100
)

func init() {
	client = request.NewClient("https://api.github.com/")
	client.SetHeaders(http.Header{
		"Accept":        []string{"application/vnd.github+json"},
		"Authorization": []string{"Token " + os.Getenv("GITHUB_OAUTH_TOKEN")},
	})
}

func GetAllRepos() ([]GHRepo, error) {
	var allRepos []GHRepo

	userRepos, err := GetCurrentUserRepos()

	if err != nil {
		return nil, err
	}

	allRepos = append(allRepos, userRepos...)

	contributionRepos, err := GetCurrentUserContributionRepos()

	if err != nil {
		return nil, err
	}

	allRepos = append(allRepos, contributionRepos...)

	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].PushedAt > allRepos[j].PushedAt
	})

	return allRepos, nil
}

func GetCurrentUserRepos() ([]GHRepo, error) {
	var page = 1
	var allRepos []GHRepo

	for {
		repos, err := GetCurrentUserReposWithPagination(page)

		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if len(repos) < queryLimit || len(repos) == 0 {
			break
		}

		page = page + 1
	}

	return allRepos, nil
}

func GetCurrentUserContributionRepos() ([]GHRepo, error) {
	res, err := query(client.Post("graphql", nil, map[string]string{"query": `query {
		viewer {
			contributionsCollection {
				pullRequestContributionsByRepository(maxRepositories: 100) {
					repository {
						isFork
						nameWithOwner
						owner {
							login
						}
						url
					}
					contributions(last: 1, orderBy:{direction:ASC}) {
						nodes {
							pullRequest {
								createdAt
							}
						}
					}
				}
			}
		}
	}`}))

	if err != nil {
		return nil, err
	}

	contributions := GHContributions{}

	if err = json.Unmarshal(res.Data, &contributions); err != nil {
		return nil, err
	}

	var repos []GHRepo

	for _, contribution := range contributions.Data.Viewer.ContributionsCollection.PullRequestContributionsByRepository {
		if contribution.Repository.Owner.Login == "HiDeoo" || contribution.Repository.IsFork || len(contribution.Contributions.Nodes) == 0 {
			continue
		}

		repos = append(repos, GHRepo{
			ID:       contribution.Repository.ID,
			FullName: contribution.Repository.NameWithOwner,
			HtmlURL:  contribution.Repository.URL,
			PushedAt: contribution.Contributions.Nodes[0].PullRequest.CreatedAt,
		})
	}

	return repos, nil
}

func GetCurrentUserReposWithPagination(page int) ([]GHRepo, error) {
	queryParams := url.Values{}
	queryParams.Set("sort", "pushed")
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("per_page", strconv.Itoa(queryLimit))

	res, err := query(client.Get("user/repos", queryParams))

	if err != nil {
		return nil, err
	}

	repos := []GHRepo{}

	if err = json.Unmarshal(res.Data, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func query(res *request.Response, err error) (*request.Response, error) {
	if res.StatusCode != 200 {
		ghError := &GHError{}

		if err := json.Unmarshal(res.Data, &ghError); err != nil || ghError.Message == "" {
			return nil, fmt.Errorf("unable to fetch GitHub data (error: %d)", res.StatusCode)
		}

		return nil, errors.New(ghError.Message)
	}

	return res, err
}
