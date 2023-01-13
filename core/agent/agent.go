package agent

import (
	"bytes"
	"chasqi-go/types"
	"encoding/json"
	"log"
	"time"
)

type Agent struct {
	idx      int
	tree     *types.Tree
	resultCh chan types.AgentResult
	visitor  NodeVisitor
	stopped  bool
}

func New(idx int,
	t *types.Tree,
	resultCh chan types.AgentResult,
	visitor NodeVisitor) *Agent {
	return &Agent{
		idx:      idx,
		tree:     t,
		resultCh: resultCh,
		visitor:  visitor,
	}
}

func (a *Agent) Start() {
	log.Printf("Agent %d started / tree: %s", a.idx, a.tree.ID)

	currentNode := a.tree.Root
	var resultSet []*types.ResponseResult
	testResult := types.AgentResult{}
	succCount := 0
	errCount := 0

	for currentNode != nil && !a.stopped {
		result, err := a.visit(currentNode)
		if err != nil {
			log.Printf("Agent %d failed to visit node %s: %s", a.idx, currentNode.Path, err)
		}

		if result.StatusCode > 299 {
			errCount++
		} else {
			succCount++
		}
		a.enrichResult(result, *currentNode)
		resultSet = append(resultSet, result)

		nextEdge := currentNode.Next
		if nextEdge != nil {
			time.Sleep(time.Duration(nextEdge.Delay) * time.Second)
			currentNode = nextEdge.After
		} else {
			// we reached the end of the tree
			currentNode = nil
		}
	}

	// Enrich the result
	testResult.SuccessCount = succCount
	testResult.ErrorCount = errCount
	testResult.Result = resultSet
	testResult.TreeID = a.tree.ID
	testResult.AgentID = a.idx
	testResult.FinishedAt = time.Now()
	a.resultCh <- testResult
}

func (a *Agent) visit(n *types.Node) (*types.ResponseResult, error) {
	// Unfortunately the nil interface is not really nil in go
	if n.Body != nil {
		var b *bytes.Buffer
		b = new(bytes.Buffer)
		json.NewEncoder(b).Encode(n.Body)
		return a.visitor.Visit(
			n.Method,
			n.Path,
			b,
			n.Headers,
		)
	} else {
		return a.visitor.Visit(
			n.Method,
			n.Path,
			nil,
			n.Headers,
		)
	}
}

func (a *Agent) enrichResult(r *types.ResponseResult, n types.Node) {
	r.Body = n.Body
	r.Method = n.Method
	r.Headers = n.Headers
}

func (a *Agent) Stop() {
	a.stopped = true
}
