package premium

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PatreonClient struct {
	httpClient         *http.Client
	proxyUrl, proxyKey string
}

const PatreonGracePeriod = time.Hour * 24 * 7

func NewPatreonClient(proxyUrl, proxyKey string) *PatreonClient {
	return &PatreonClient{
		httpClient: &http.Client{
			Timeout: time.Second * 1,
		},
		proxyUrl: proxyUrl,
		proxyKey: proxyKey,
	}
}

func NewPatreonClientWithHttpClient(httpClient *http.Client, proxyUrl, proxyKey string) *PatreonClient {
	return &PatreonClient{
		httpClient: httpClient,
		proxyUrl:   proxyUrl,
		proxyKey:   proxyKey,
	}
}

type proxyResponse struct {
	Premium bool
	Tier    int
}

func (p *PatreonClient) GetTier(ctx context.Context, userIds ...uint64) (PremiumTier, error) {
	strIds := make([]string, len(userIds))
	for i, userId := range userIds {
		strIds[i] = strconv.FormatUint(userId, 10)
	}

	url := fmt.Sprintf("%s/ispremium?id=%s", p.proxyUrl, strings.Join(strIds, ","))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return None, err
	}

	req.Header.Set("Authorization", p.proxyKey)

	res, err := p.httpClient.Do(req)
	if err != nil {
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
