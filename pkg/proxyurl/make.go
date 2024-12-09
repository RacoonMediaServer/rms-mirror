package proxyurl

import (
	"encoding/base64"
	"net/url"
)

type Maker struct {
	BaseURL *url.URL
}

func (m Maker) MakeURL(contentURL string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(contentURL))
	u := *m.BaseURL

	query := u.Query()
	query.Add("href", encoded)

	u.RawQuery = query.Encode()
	return u.String()
}
