package webproxy

import (
	"net/http"
	"net/url"
)

type WebProxy struct {
	client              *http.Client
	proxyUrl            string
	authHeader, authKey string
}

func NewWebProxy(proxyUrl, authHeader, authKey string) *WebProxy {
	return &WebProxy{
		client:     &http.Client{},
		proxyUrl:   proxyUrl,
		authHeader: authHeader,
		authKey:    authKey,
	}
}

func (p *WebProxy) Do(req *http.Request) (*http.Response, error) {
	dest := req.URL.String()

	var err error
	req.URL, err = url.Parse(p.proxyUrl)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("url", dest)
	req.URL.RawQuery = query.Encode()

	req.Host = req.URL.Host

	req.Header.Set(p.authHeader, p.authKey)

	return p.client.Do(req)
}
