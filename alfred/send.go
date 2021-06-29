package alfred

import (
	"encoding/json"
	"fmt"
	"log"
)

type Result struct {
	Items []interface{} `json:"items"`
}

type BaseItem struct {
	Title    string `json:"title"`
	SubTitle string `json:"subtitle"`
}

type Item struct {
	BaseItem
	Arg string `json:"arg"`
}

type EmptyItem struct {
	Item
	Mods Modifiers `json:"mods"`
}

type Modifiers struct {
	Alt   Modifier `json:"alt"`
	Cmd   Modifier `json:"cmd"`
	Shift Modifier `json:"shift"`
}

type Modifier struct {
	Valid bool `json:"valid"`
}

type Error struct {
	BaseItem
	Valid bool `json:"valid"`
}

func SendResult(items []Item) {
	var result Result

	if len(items) > 0 {
		result = Result{Items: make([]interface{}, len(items))}

		for i, item := range items {
			result.Items[i] = item
		}
	} else {
		result = Result{Items: []interface{}{newEmptyPlaceholderItem()}}
	}

	send(result)
}

func SendError(err error) {
	send(Result{
		Items: []interface{}{
			Error{
				BaseItem: BaseItem{"Something went wrong!", err.Error()},
				Valid:    false,
			},
		},
	})
}

func send(data interface{}) {
	bytes, err := json.Marshal(data)

	if err == nil {
		fmt.Println(string(bytes))

		return
	}

	log.Panicln(err)
}

func newEmptyPlaceholderItem() EmptyItem {
	return EmptyItem{
		Item: Item{
			BaseItem: BaseItem{"You're alone! ¯\\_(ツ)_/¯", "Try browsing Twitch…"},
			Arg:      "https://www.twitch.tv/directory/following",
		},
		Mods: Modifiers{
			Alt:   Modifier{false},
			Cmd:   Modifier{false},
			Shift: Modifier{false},
		},
	}
}
