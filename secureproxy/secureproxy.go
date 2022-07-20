package secureproxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TicketsBot/common/sentry"
	"github.com/TicketsBot/common/utils"
	"io/ioutil"
	"net/http"
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
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    *string           `json:"body,omitempty"`
}

func (p *Client) DoRequest(method, url string, headers map[string]string, bodyData []byte) ([]byte, error) {
	body := secureProxyRequest{
		Method:  method,
		Url:     url,
		Headers: headers,
	}

	if bodyData != nil {
		body.Body = utils.Ptr(string(bodyData))
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Post(p.Url+"/proxy", "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		sentry.Error(err)
		return nil, errors.New("error proxying request")
	}

	defer res.Body.Close()

	if errorHeader := res.Header.Get("x-proxy-error"); errorHeader != "" {
		return nil, errors.New(errorHeader)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("integration request returned status code %d", res.StatusCode)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
