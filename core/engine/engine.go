package engine

import (
	"chasqi-go/core/agent"
	"chasqi-go/types"
	"fmt"
	"log"
	"sync"
	"time"
)

type DefaultEngine struct {
	statusMap      map[types.TreeID]*types.LoopStatus
	activeTrees    map[types.TreeID]*types.Tree
	doneTrees      map[types.TreeID]*types.Tree
	resultMap      map[types.TreeID]*types.TestResult
	enqueuedTrees  []*types.Tree
	visitorCreator func() agent.NodeVisitor
	resultCh       chan types.TestResult
	mu             *sync.Mutex
	exitCh         chan struct{}
	hasStopped     bool
}

func (e *DefaultEngine) ById(id string) *types.LoopStatus {

	return e.statusMap[types.TreeID(id)]
}

func (e *DefaultEngine) Enqueue(tree *types.Tree) error {
	if e.hasStopped {
		return fmt.Errorf("DefaultEngine has stopped")
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	e.enqueuedTrees = append(e.enqueuedTrees, tree)
	return nil
}

func (e *DefaultEngine) Start() {
	log.Printf("DefaultEngine started")
	loopTimer := time.NewTicker(3 * time.Second)
	flushTimer := time.NewTicker(5 * 60 * time.Second)
coreLoop:
	for {
		select {
		case <-loopTimer.C:
			e.onTick()
		case <-flushTimer.C:
			break
		case result := <-e.resultCh:
			log.Printf("got result with %d entries", len(result.Result))
			e.onResult(&result)
		case <-e.exitCh:
			break coreLoop
		}
	}
}

func (e *DefaultEngine) onTick() {
	e.mu.Lock()
	for _, t := range e.enqueuedTrees {
		log.Printf("processing tree: %v", t)
		newEnqueuedTrees := make([]*types.Tree, 0)

		for _, t2 := range e.enqueuedTrees {
			if t2 != t {
				newEnqueuedTrees = append(newEnqueuedTrees, t2)
			}
		}

		s := time.Now()
		e.statusMap[types.TreeID(t.ID)] = &types.LoopStatus{
			TreeID:    t.ID,
			IsDone:    false,
			StartedAt: &s,
		}
		e.visitTree(t)
		e.enqueuedTrees = newEnqueuedTrees
	}
	e.mu.Unlock()
}

func (e *DefaultEngine) onResult(result *types.TestResult) {
	e.mu.Lock()
	defer e.mu.Unlock()
	t := e.activeTrees[types.TreeID(result.TreeID)]
	delete(e.activeTrees, types.TreeID(result.TreeID))
	e.doneTrees[types.TreeID(result.TreeID)] = t
	e.statusMap[types.TreeID(result.TreeID)].IsDone = true
	e.resultMap[types.TreeID(result.TreeID)] = result

	log.Printf("Finished tree: %s", result.String())
}

func (e *DefaultEngine) visitTree(tree *types.Tree) {
	n := tree.Config.AgentAmount
	for i := 0; i < n; i++ {
		go func(i int) {
			a := agent.New(
				i,
				tree,
				e.resultCh,
				e.visitorCreator())
			a.Start()
		}(i)
	}
}

func (e *DefaultEngine) Get(id string) *types.TestResult {
	return e.resultMap[types.TreeID(id)]
}

func (e *DefaultEngine) Cancel(id string) error {
	//TODO implement me
	panic("implement me")
}

func New(visitorCreator func() agent.NodeVisitor, exitCh chan struct{}) *DefaultEngine {
	return &DefaultEngine{
		statusMap:      make(map[types.TreeID]*types.LoopStatus),
		activeTrees:    make(map[types.TreeID]*types.Tree),
		doneTrees:      make(map[types.TreeID]*types.Tree),
		resultMap:      make(map[types.TreeID]*types.TestResult),
		enqueuedTrees:  make([]*types.Tree, 0),
		visitorCreator: visitorCreator,
		resultCh:       make(chan types.TestResult, 1000),
		exitCh:         exitCh,
		mu:             &sync.Mutex{},
	}
}
