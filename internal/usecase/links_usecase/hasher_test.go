package links_usecase

import "testing"

func TestHasher(t *testing.T) {
	inputs := []struct {
		testName    string
		originalUrl string
		shortUrl    string
	}{
		{
			testName:    "first good scenario",
			originalUrl: "hashMe",
			shortUrl:    "8ZUzT7sqbX",
		},
		{
			testName:    "second good scenario",
			originalUrl: "hochu/v/Ozon!",
			shortUrl:    "kxmpOd9sLi",
		},
		{
			testName:    "third good scenario",
			originalUrl: "!! catch/ми/if/ю/can !!",
			shortUrl:    "094kJp5rbh",
		},
	}

	for _, tt := range inputs {
		t.Run(tt.testName, func(t *testing.T) {
			if shortUrl := Hasher(tt.originalUrl); shortUrl != tt.shortUrl {
				t.Errorf("want %s, got %s", tt.shortUrl, shortUrl)
			}
		})
	}
}
