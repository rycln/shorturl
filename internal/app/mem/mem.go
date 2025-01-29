package mem

type Memorizer interface {
	AddURL(string)
	GetURL(string) string
}

type MemStorage struct {
	Storage map[string]string
	incId   int64
}

func (m MemStorage) AddURL(url string) {

}

func (m MemStorage) GetURL(shortURL string) string {

}
