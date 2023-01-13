package visitor

import (
	"chasqi-go/types"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (c *defaultHttpClient) Visit(method, url string, body io.Reader, headers map[string][]string) (*types.ResponseResult, error) {
	s := time.Now()

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("error while creating request: %s", err)
		return nil, err
	}
	for k, v := range headers {
		req.Header[k] = v
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
