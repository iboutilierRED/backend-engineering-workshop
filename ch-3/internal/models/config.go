package models

type AppConfig struct {
	FeedUrl   string `json:"feed_url"`
	StatsPort string `json:"stats_port"`
}
