package delivery

type CreateShortUrlParams struct {
	OriginalUrl string `json:"original_url"`
}

type CreateShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type GetOriginalUrlResp struct {
	OriginalUrl string `json:"original_url"`
}
