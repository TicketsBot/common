package secureproxy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TicketsBot/common/sentry"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Client struct {
	Url    string
	client *http.Client
}

func NewSecureProxy(url string) *Client {
	return &Client{
		Url:    url,
		client: &http.Client{},
	}
}

type secureProxyRequest struct {
	Method   string            `json:"method"`
	Url      string            `json:"url"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     []byte            `json:"body,omitempty"`
	JsonBody json.RawMessage   `json:"json_body,omitempty"`
}

type requestBody interface {
	[]byte | any
}

func (p *Client) DoRequest(method, url string, headers map[string]string, bodyData requestBody) ([]byte, int, error) {
	body := secureProxyRequest{
		Method:  method,
		Url:     url,
		Headers: headers,
	}

	// nil will fall through
	switch v := bodyData.(type) {
	case []byte:
		base64.StdEncoding.Encode(body.Body, v)
	case any:
		encoded, err := json.Marshal(v)
		if err != nil {
			return nil, 0, err
		}

		body.JsonBody = json.RawMessage(encoded)
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, 0, err
	}

	res, err := p.client.Post(p.Url+"/proxy", "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		sentry.Error(err)
		return nil, 0, errors.New("error proxying request")
	}

	defer res.Body.Close()

	if errorHeader := res.Header.Get("x-proxy-error"); errorHeader != "" {
		return nil, 0, errors.New(errorHeader)
	}

	if res.StatusCode != 200 {
		return nil, res.StatusCode, fmt.Errorf("integration request returned status code %d", res.StatusCode)
	}

	statusRaw := res.Header.Get("x-status-code")
	if statusRaw == "" {
		return nil, 0, errors.New("response missing x-status-code header")
	}

	statusCode, err := strconv.Atoi(statusRaw)
	if err != nil {
		return nil, 0, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return resBody, statusCode, nil
}
