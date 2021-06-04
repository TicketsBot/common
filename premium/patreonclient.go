package premium

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func (p *PatreonClient) GetTier(userIds ...uint64) (PremiumTier, error) {
	strIds := make([]string, len(userIds))
	for i, userId := range userIds {
		strIds[i] = strconv.FormatUint(userId, 10)
	}

	url := fmt.Sprintf("%s/ispremium?key=%s&id=%s", p.proxyUrl, p.proxyKey, strings.Join(strIds, ","))
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
