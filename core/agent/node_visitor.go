package agent

import "chasqi-go/types"

type NodeVisitor interface {
	Get(host, url string, headers map[string][]string) (*types.ResponseResult, error)
	Post(host, url string, headers map[string][]string, body map[string][]string) (*types.ResponseResult, error)
	Put(host, url string, headers map[string][]string, body map[string][]string) (*types.ResponseResult, error)
	Delete(host, url string, headers map[string][]string, body map[string][]string) (*types.ResponseResult, error)
}
