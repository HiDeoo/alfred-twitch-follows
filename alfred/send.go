package alfred

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func SendResult(items []Item) {
	send(Result{Items: items})
}

func SendError(err error) {
	send(Result{
		Items: []interface{}{
			Error{
				Title:    "Something went wrong!",
				SubTitle: err.Error(),
				Valid:    false,
			},
		},
	})
}

func send(data interface{}) {
	bytes, err := json.Marshal(data)

	if err == nil {
		fmt.Println(string(bytes))

		os.Exit(1)
	}

	log.Panicln(err)
}

type Result struct {
	Items interface{} `json:"items"`
}

type Item struct {
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
}

type Error struct {
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
	Valid    bool   `json:"valid"`
}
