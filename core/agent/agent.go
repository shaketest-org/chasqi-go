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
	// currentNode Node that should be visited
	currentNode := a.tree.Root
	log.Printf("Agent %d started / tree: %s", a.idx, a.tree.ID)
	var resultSet []*types.ResponseResult

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
		resultSet = append(resultSet, result)
		nextEdge := currentNode.Next

		if nextEdge != nil {
			time.Sleep(time.Duration(nextEdge.Delay) * time.Second)
			currentNode = nextEdge.After
		} else {
			currentNode = nil
		}
	}

}

func (a *Agent) visit() {

}

func (a *Agent) Stop() {
	a.stopped = true
}
