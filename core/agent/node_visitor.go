package agent

import (
	"chasqi-go/types"
	"io"
)

type NodeVisitor interface {
	Visit(method, url string, body io.Reader, headers map[string][]string) (*types.ResponseResult, error)
}
