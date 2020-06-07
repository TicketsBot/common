package premium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PatreonClient struct {
	httpClient *http.Client
	proxyUrl, proxyKey string
}

func NewPatreonClient(proxyUrl, proxyKey string) *PatreonClient {
	return &PatreonClient{
		httpClient: &http.Client{
			Timeout: time.Second * 3,
		},
		proxyUrl: proxyUrl,
		proxyKey: proxyKey,
	}
}

type proxyResponse struct {
	Premium bool
	Tier int
}

func (p *PatreonClient) GetTier(userId uint64) (PremiumTier, error) {
	url := fmt.Sprintf("%s/ispremium?key=%s&id=%d", p.proxyUrl, p.proxyKey, userId)
	res, err := p.httpClient.Get(url); if err != nil {
		return None, err
	}

	var data proxyResponse
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return None, err
	}

	if data.Premium {
		return PremiumTier(data.Tier), nil
	}

	return None, nil
}
