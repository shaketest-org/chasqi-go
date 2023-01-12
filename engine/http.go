package engine

import (
	"chasqi-go/types"
	"fmt"
	"net/http"
	url2 "net/url"
	"time"
)

type defaultHttpClient struct {
	Client *http.Client
}

func NewDefaultHttpClient() *defaultHttpClient {
	return &defaultHttpClient{
		Client: &http.Client{},
	}
}

func (c *defaultHttpClient) Get(host, url string, headers map[string][]string) (*types.ResponseResult, error) {
	s := time.Now()
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url2.URL{
			RawPath: url,
		},
		Header: headers,
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while sending GET request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	e := time.Now()

	ms := e.Sub(s).Milliseconds()

	return &types.ResponseResult{
		StatusCode:   resp.StatusCode,
		DurationInMs: ms,
	}, nil
}
