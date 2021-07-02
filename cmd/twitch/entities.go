package main

type TwitchFollows struct {
	Total      int              `json:"total"`
	Data       []TwitchFollow   `json:"data"`
	Pagination TwitchPagination `json:"pagination"`
}

type TwitchFollow struct {
	FromId     string `json:"from_id"`
	FromLogin  string `json:"from_login"`
	ToId       string `json:"to_id"`
	ToLogin    string `json:"to_login"`
	ToName     string `json:"to_name"`
	FollowedAt string `json:"followed_at"`
}

type TwitchStreams struct {
	Data       []TwitchStream   `json:"data"`
	Pagination TwitchPagination `json:"pagination"`
}

type TwitchStream struct {
	Id           string   `json:"id"`
	UserId       string   `json:"user_id"`
	UserLogin    string   `json:"user_login"`
	UserName     string   `json:"user_name"`
	GameId       string   `json:"game_id"`
	GameName     string   `json:"game_name"`
	Type         string   `json:"type"`
	Title        string   `json:"title"`
	ViewerCount  int      `json:"viewer_count"`
	StartedAt    string   `json:"started_at"`
	Language     string   `json:"language"`
	ThumbnailUrl string   `json:"thumbnail_url"`
	TagIds       []string `json:"tag_ids"`
	IsMature     bool     `json:"is_mature"`
}

type TwitchUsers struct {
	Data []TwitchUser `json:"data"`
}

type TwitchUser struct {
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

type TwitchError struct {
	Error   string `json:"error"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type TwitchPagination struct {
	Cursor string `json:"cursor"`
}
