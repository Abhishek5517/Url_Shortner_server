package models

import "time"

type UserMapped struct {
	Email     string
	ShortURLs []map[string]string
}

type UrlMap struct {
	Short_code string    `json:"shortCode"`
	Actual_url string    `json:"actualUrl"`
	Hits       int64     `json:"hits"`
	CreatedAt  time.Time `json:"createdAt"`
}
