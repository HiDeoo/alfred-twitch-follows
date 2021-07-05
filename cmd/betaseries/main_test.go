package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnwatchedShowItems(t *testing.T) {
	tests := []struct {
		name      string
		showCount int
	}{
		{"ReturnNoShows", 0},
		{"ReturnSingleShow", 1},
		{"ReturnMultipleShows", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var shows = []BSShow{}

			for i := 0; i < test.showCount; i++ {
				showJSON := BSShow{
					ID:    123456789 + i,
					Title: fmt.Sprintf("Title %d", i),
					User:  BSShowUser{Next: BSShowUserNext{ID: 987654321 + 0}},
				}

				shows = append(shows, showJSON)
			}

			items, _ := getUnwatchedShowItems(func() ([]BSShow, error) {
				return shows, nil
			})

			assert.Equal(t, test.showCount, len(items))

			for i, item := range items {
				assert.Equal(t, shows[i].Title, item.Title)
				assert.Equal(
					t,
					fmt.Sprintf("%d episodes remaining (%s total)", shows[i].User.Remaining, shows[i].Episodes),
					item.SubTitle,
				)
				assert.Equal(t, strconv.Itoa(shows[i].User.Next.ID), item.Arg)
			}
		})
	}
}
