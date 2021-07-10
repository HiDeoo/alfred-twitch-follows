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
	Arg  string     `json:"arg"`
	Mods *Modifiers `json:"mods,omitempty"`
}

type EmptyItem struct {
	Item
	Icon EmptyItemIcon `json:"icon"`
	Mods Modifiers     `json:"mods"`
}

type EmptyItemIcon struct {
	Path string `json:"path"`
}

type Modifiers struct {
	Alt   Modifier `json:"alt"`
	Cmd   Modifier `json:"cmd"`
	Shift Modifier `json:"shift"`
}

type Modifier struct {
	Valid    bool   `json:"valid"`
	Arg      string `json:"arg,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

type Error struct {
	BaseItem
	Valid bool          `json:"valid"`
	Icon  EmptyItemIcon `json:"icon"`
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
				Icon:     EmptyItemIcon{"images/error.png"},
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
		Icon: EmptyItemIcon{"images/error.png"},
		Mods: Modifiers{
			Alt:   Modifier{Valid: false},
			Cmd:   Modifier{Valid: false},
			Shift: Modifier{Valid: false},
		},
	}
}
