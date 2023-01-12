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
	resultCh chan types.TestResult
	visitor  NodeVisitor
	stopped  bool
}

func New(idx int,
	t *types.Tree,
	resultCh chan types.TestResult,
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
	testResult := types.TestResult{}
	succCount := 0
	errCount := 0
	s := time.Now()
	testResult.StartedAt = &s

	for currentNode != nil && !a.stopped {
		var b *bytes.Buffer
		if currentNode.Body != nil {
			b = new(bytes.Buffer)
			json.NewEncoder(b).Encode(currentNode.Body)
		}
		result, err := a.visitor.Visit(
			currentNode.Method,
			currentNode.Path,
			b,
			currentNode.Headers,
		)
		if err != nil {
			log.Printf("Agent %d failed to visit node %s: %s", a.idx, currentNode.Path, err)
		}

		if result.StatusCode > 299 {
			errCount++
		} else {
			succCount++
		}

		result.Node = currentNode
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

	e := time.Now()
	testResult.FinishedAt = &e
	testResult.SuccessCount = succCount
	testResult.ErrorCount = errCount
	testResult.Result = resultSet
	testResult.TreeID = a.tree.ID

	a.resultCh <- testResult

}

func (a *Agent) Stop() {
	a.stopped = true
}
