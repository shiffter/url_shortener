package delivery

type CreateShortUrlRequest struct {
	OriginalUrl string `json:"original_url"`
}

type CreateShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
	Status   int    `json:"status"`
	Error    string
}

type GetOriginalUrlResponse struct {
	OriginalUrl string `json:"original_url"`
	Status      int    `json:"status"`
	Error       string
}
