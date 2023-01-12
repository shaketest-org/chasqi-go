package agent

import (
	"chasqi-go/types"
	"log"
)

type Agent struct {
	idx      int
	tree     *types.Tree
	resultCh chan types.TestResult
	visitor  NodeVisitor
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
	// current Node that should be visited
	c := a.tree.Root
	log.Printf("Agent %d started / root: %s", a.idx, c.String())
	// Next Edge that points to the next Node
	e := c.Next.After
	for e != nil {

	}

}

func (a *Agent) visit() {

}

func (a *Agent) Stop() {
}
