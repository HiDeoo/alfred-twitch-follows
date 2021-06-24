package alfred

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func SendError(err error) {
	result := Result{
		Items: []interface{}{
			Error{
				Title:    "Something went wrong!",
				SubTitle: err.Error(),
				Valid:    false,
			},
		},
	}

	bytes, err := json.Marshal(result)

	if err == nil {
		fmt.Println(string(bytes))

		os.Exit(1)
	}

	log.Panicln(err)
}

type Result struct {
	Items []interface{} `json:"items"`
}

type Error struct {
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
	Valid    bool   `json:"valid"`
}
