package storage

type ShortenedURL struct {
	ShortURL string `json:"short_url"`
	OrigURL  string `json:"original_url"`
}

func NewShortenedURL(shortURL, origURL string) ShortenedURL {
	surl := ShortenedURL{
		ShortURL: shortURL,
		OrigURL:  origURL,
	}
	return surl
}
