package twitch

type Follows struct {
	Total      int        `json:"total"`
	Data       []Follow   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Follow struct {
	FromId     string `json:"from_id"`
	FromLogin  string `json:"from_login"`
	ToId       string `json:"to_id"`
	ToLogin    string `json:"to_login"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

type Users struct {
	Data []User `json:"data"`
}

type User struct {
	Id              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	CreatedAt       string `json:"created_at"`
}

type Error struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type QueryParam struct {
	Key   string
	Value string
}
