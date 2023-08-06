package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepos(t *testing.T) {
	tests := []struct {
		name      string
		repoCount int
	}{
		{"ReturnNoRepos", 0},
		{"ReturnSingleRepo", 1},
		{"ReturnMultipleRepos", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var repos = []GHRepo{}

			for i := 0; i < test.repoCount; i++ {
				repoJSON := GHRepo{
					ID:       123456789 + i,
					FullName: fmt.Sprintf("Full name %d", i),
					HtmlURL:  fmt.Sprintf("https://github.com/user/repo%d", i),
				}

				repos = append(repos, repoJSON)
			}

			items, _ := getRepos(func() ([]GHRepo, error) {
				return repos, nil
			})

			assert.Equal(t, test.repoCount, len(items))

			for i, item := range items {
				assert.Equal(t, repos[i].FullName, item.Title)
				assert.Regexp(t, regexp.MustCompile("^https://github.com/"), item.SubTitle)
				assert.Equal(t, repos[i].HtmlURL, item.Arg)
			}
		})
	}
}
