package agent

import "chasqi-go/types"

type NodeVisitor interface {
	Get(url string, headers map[string][]string) (*types.ResponseResult, error)
	Post(url string, headers map[string][]string, body map[string][]string) (*types.ResponseResult, error)
	Put(url string, headers map[string][]string, body map[string][]string) (*types.ResponseResult, error)
	Delete(url string, headers map[string][]string, body map[string][]string) (*types.ResponseResult, error)
}
