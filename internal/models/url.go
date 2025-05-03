package models

type ShortURL string

type OrigURL string

type URLPair struct {
	UID   UserID
	Short ShortURL
	Orig  OrigURL
}
