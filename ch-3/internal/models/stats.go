package models

type Stats struct {
	Messages    int `json:"messages_consumed"`
	UniqueUsers int `json:"distinct_users"`
	UniqueUris  int `json:"distinct_uris"`
	Bots        int `json:"num_bots"`
	NonBots     int `json:"num_non_bots"`
}
