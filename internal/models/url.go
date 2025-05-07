package models

type ShortURL string

type OrigURL string

type URLPair struct {
	UID   UserID   `json:"user_id"`
	Short ShortURL `json:"short_url"`
	Orig  OrigURL  `json:"original_url"`
}

type DelURLReq struct {
	UID   UserID   `json:"user_id"`
	Short ShortURL `json:"short_url"`
}
